#ifndef __TYPES_H__
#define __TYPES_H__


#include <stdint.h>

typedef uint64_t KeyType;
typedef uint32_t ValType;

typedef struct {
	KeyType key;
	ValType val;
} KVPair;

#endif
