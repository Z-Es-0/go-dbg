#include "windows.h"
#include "stdio.h"

const char szSecNameText[8]    = {'.', 't', 'e', 'x', 't', 0, 0, 0 };
const char szSecNameTextBSS[8] = {'.', 't', 'e', 'x', 't', 'b', 's', 's'};
const char szSecNameTLS[8]     = {'.', 't', 'l', 's', 0, 0, 0};
const char szSecNameBSS[8]     = {'.', 'b', 's', 's', 0, 0, 0, 0};
const char szSecNameData[8]    = {'.', 'd', 'a', 't', 'a', 0, 0, 0};


BOOL _IsPack_32(PBYTE pBuf, UINT uSize)
{
	BOOL bRet = FALSE;

	PIMAGE_DOS_HEADER      pDosHdr  = NULL;
	PIMAGE_NT_HEADERS      pNTHdr   = NULL;
	PIMAGE_FILE_HEADER     pFileHdr = NULL;
	PIMAGE_SECTION_HEADER  pSecHdr  = NULL;
	UINT				   i        = 0;
	
	pDosHdr  = (PIMAGE_DOS_HEADER)pBuf;
	pNTHdr   = (PIMAGE_NT_HEADERS)((PBYTE)pDosHdr + pDosHdr->e_lfanew);
	pFileHdr = &pNTHdr->FileHeader;
	pSecHdr  = (PIMAGE_SECTION_HEADER)((PBYTE)pNTHdr + sizeof(IMAGE_NT_HEADERS32));

	for(i = 0; i < pFileHdr->NumberOfSections; i++)
	{
		if (IsBadReadPtr(pSecHdr, sizeof(IMAGE_SECTION_HEADER)))
			break;

		// Rule1
		if (pSecHdr[i].Characteristics & IMAGE_SCN_CNT_CODE && 0x1000 < pSecHdr[i].Misc.VirtualSize && pSecHdr[i].SizeOfRawData <= pSecHdr[i].Misc.VirtualSize - 0x1000)
			goto Exit1;

		// Rule2
		if (!memcmp(pSecHdr[i].Name, szSecNameText, sizeof(szSecNameText)) && 0x1000 < pSecHdr[i].Misc.VirtualSize && pSecHdr[i].SizeOfRawData < pSecHdr[i].Misc.VirtualSize - 0x1000)
			goto Exit1;

		// Rule3
		if (0 == *pSecHdr[i].Name && pSecHdr[i].SizeOfRawData < pSecHdr[i].Misc.VirtualSize)
			goto Exit1;

		// Rule4
		if (0 != pSecHdr[i].SizeOfRawData)
			continue;

		// Rule5
		if (0x1000 > pSecHdr[i].Misc.VirtualSize)
			continue;

		// Rule6
		if (!memcmp(pSecHdr[i].Name, szSecNameTextBSS, sizeof(szSecNameTextBSS)))
			continue;

		// Rule7
		if (!memcmp(pSecHdr[i].Name, szSecNameTLS, sizeof(szSecNameTLS)))
			continue;
		
		// Rule8
		if (!memcmp(pSecHdr[i].Name, szSecNameBSS, sizeof(szSecNameBSS)))
			continue;

		// Rule9
		if (!memcmp(pSecHdr[i].Name, szSecNameData, sizeof(szSecNameData)))
			continue;

		// Rule10
		if (pSecHdr[i].Characteristics & (IMAGE_SCN_CNT_INITIALIZED_DATA | IMAGE_SCN_CNT_UNINITIALIZED_DATA))
			goto Exit1;
	}
	goto Exit0;
Exit1:
	bRet = TRUE;
Exit0:
	return bRet;
}


BOOL _IsPack_64(PBYTE pBuf, UINT uSize)
{
	// no support
	return FALSE;
}

BOOL IsPack(PBYTE pBuf, UINT uSize)
{
	BOOL bRet = FALSE;

	PIMAGE_DOS_HEADER   pDosHdr  = NULL;
	PIMAGE_NT_HEADERS   pNTHdr   = NULL;
	PIMAGE_FILE_HEADER  pFileHdr = NULL;
 
	pDosHdr  = (PIMAGE_DOS_HEADER)pBuf;
	pNTHdr   = (PIMAGE_NT_HEADERS)((PBYTE)pDosHdr + pDosHdr->e_lfanew);
	
	pFileHdr = &pNTHdr->FileHeader;

	if (IsBadReadPtr(pFileHdr, sizeof(IMAGE_FILE_HEADER)))
		goto Exit0;

	if (!(pFileHdr->Characteristics |  IMAGE_FILE_32BIT_MACHINE) || pFileHdr->Machine == IMAGE_FILE_MACHINE_AMD64 || pFileHdr->Machine == IMAGE_FILE_MACHINE_IA64)
		bRet = _IsPack_64(pBuf, uSize);
	else
		bRet = _IsPack_32(pBuf, uSize);

Exit0:
	return bRet;	

}


BOOL IsPE(PBYTE pBuf, UINT uSize)
{
	BOOL bRet = FALSE;

	PIMAGE_DOS_HEADER   pDosHdr  = NULL;
	PIMAGE_NT_HEADERS   pNTHdr   = NULL;


	if (NULL == pBuf || 0 == uSize)
		goto Exit0;

	
	pDosHdr = (PIMAGE_DOS_HEADER)pBuf;
	if (IsBadReadPtr(pDosHdr, sizeof(IMAGE_DOS_HEADER)))
		goto Exit0;
	if (IMAGE_DOS_SIGNATURE != pDosHdr->e_magic)
		goto Exit0;

	pNTHdr = (PIMAGE_NT_HEADERS)((PBYTE)pDosHdr + pDosHdr->e_lfanew);
	if (IsBadReadPtr(pNTHdr, sizeof(IMAGE_NT_HEADERS32)))
		goto Exit0;
	if (IMAGE_NT_SIGNATURE != pNTHdr->Signature)
		goto Exit0;

	bRet = TRUE;
Exit0:
	return bRet;
}

int main(int argc, char *argv[])
{

	DWORD	dwAttrib = 0; 
	HANDLE	hFile    = NULL;
	HANDLE  hMapping = NULL;
	PBYTE	pBuf     = NULL;
	UINT	uSize    = 0; 


	if (argc < 2)
	{
		printf("I need a param!\n");
		return 1;
	}

	PCHAR pTest = "c:\\windows\\notepad.exe";



	dwAttrib = GetFileAttributes(argv[1]);



	if (INVALID_FILE_ATTRIBUTES == dwAttrib || FILE_ATTRIBUTE_DIRECTORY == dwAttrib)
	{
		printf("'%s' is not a file\n", argv[1]);
		return 1;
	}

	hFile = CreateFile((LPCTSTR)argv[1], GENERIC_READ, FILE_SHARE_READ, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
	if (INVALID_HANDLE_VALUE == hFile)
	{
		printf("'%s' can not open\n", argv[1]);
		return 1;
	}


	hMapping = CreateFileMapping(hFile, NULL, PAGE_READONLY, 0, 0, NULL);
	if (INVALID_HANDLE_VALUE == hMapping)
	{
		printf("'%s' can not map\n", argv[1]);
		return 1;
	}

	pBuf = (PBYTE)MapViewOfFile(hMapping, FILE_MAP_READ, 0, 0, 0);
	if (NULL == pBuf)
	{
		printf("'%s' can not MapViewOfFile\n", argv[1]);
		return 1;
	}

	uSize = GetFileSize(hFile, 0);
	if (0 == uSize)
	{
		printf("'%s' size is 0!\n", argv[1]);
		return 1;
	}

	if (!IsPE(pBuf, uSize))
	{
		printf("This file is not a PE file!\n");
		return 1;
	}

	if (IsPack(pBuf, uSize))
	{
		printf("This PE file is packed!\n");
	}
	else
	{
		printf("This PE file is NOT packed!\n");
	}

if (hFile)
{
	CloseHandle(hFile);
	hFile = NULL;
}

if (hMapping)
{
	CloseHandle(hMapping);
	hMapping = NULL;
}

if (pBuf)
{
	UnmapViewOfFile(pBuf);
	pBuf = NULL;
}

return 0;
}
