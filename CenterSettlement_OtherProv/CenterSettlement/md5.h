#ifndef BZF_MD5_H  
#define BZF_MD5_H  

#pragma once

//typedef unsigned long int UINT32;
//typedef unsigned short int UINT16;

/* MD5 context. */
typedef struct {
	UINT32 state[4];                                   /* state (ABCD) */
	UINT32 count[2];        /* number of bits, modulo 2^64 (lsb first) */
	unsigned char buffer[64];                         /* input buffer */
} MD5_CTX;

void MD5Init   (MD5_CTX *);
void MD5Update (MD5_CTX *, unsigned char *, unsigned int);
void MD5Final  (unsigned char [16], MD5_CTX *);


#endif 