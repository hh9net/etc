#include <stdio.h>
#include <dlfcn.h>

int testcode_c()
{
    void *handle;
    typedef int (*FPTR)();
    handle = dlopen("./lz77.so", 1);
    if(handle == 0){
        printf("load library lz77.so failed\n");
        return -1;
    }

    FPTR fptr = (FPTR)dlsym(handle, "Testcode");

    int iResult = (*fptr)();
    dlclose(handle);
    return iResult;
}

int compress(char *srcfile ,char *destfile)
{
    void *handle;
    typedef int (*FPTR)(char *,char *);
    handle = dlopen("./lz77.so", 1);
    if(handle == 0) {
        printf("load library lz77.so failed\n");
        return -1;
    }

    FPTR fptr = (FPTR)dlsym(handle, "Compressfile");
    int iResult = (*fptr)(srcfile, destfile);
    printf("lz77 file from %s  to %s\n", srcfile, destfile);
    return iResult;
}

int decompress(char *srcfile ,char *destfile)
{
     void *handle;
    typedef int (*FPTR)(char *,char *);
    handle = dlopen("./lz77.so", 1);
    if(handle == 0) {
        printf("load library lz77.so failed\n");
        return -1;
    }

    FPTR fptr = (FPTR)dlsym(handle, "Decompressfile");
    int iResult = (*fptr)(srcfile, destfile);
    printf("unlz77 file from %s  to %s\n", srcfile, destfile);
    return iResult;
}
