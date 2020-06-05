#ifndef __READERINF_H__
#define __READERINF_H__

#ifdef __cplusplus
extern "C" {
#endif

int Compressfilefunc(char *inputfile, char *outputfile);
int Decompressfilefunc(char *inputfile, char *outputfile);
int Testcode();

#ifdef __cplusplus
}
#endif

#endif
