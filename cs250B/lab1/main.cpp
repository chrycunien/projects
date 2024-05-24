#include <stdio.h>
#include <stdlib.h>

#include <cmath>
#include <ctime>
#include <chrono>
#include <algorithm>

#include "types.h"

extern void user_sort(KVPair* array, size_t data_cnt, size_t threads);
extern void quick_sort(KVPair* array, size_t data_cnt);
extern bool my_comparator(KVPair a, KVPair b);

int main(int argc, char** argv) {
	int threads = 8;
	int sorterchoice = 0;
	size_t data_cnt = 1024*1024*64;
	// size_t data_cnt = 32;
	// size_t data_cnt = 64000;

	if ( argc >= 2) {
		char as = argv[1][0];
		switch (as) {
			case 'Q': case 'q':{
				sorterchoice = 0;
				printf( "Choosing quicksort\n" );
				break;
			}
			case 'S': case 's':{
				sorterchoice = 1;
				printf( "Choosing std::sort\n" );
				break;
			}
			case 'U': case 'u':{
				sorterchoice = 2;
				printf( "Choosing user sort\n" );
				break;
			}
			default: {
				printf( "Choosing quicksort\n" );
				break;
			}
		}
	} else {
		printf( "Choosing quicksort\n" );
	}

	if ( argc >= 3 ) {
		int cnt = atoi(argv[2]);
		if ( cnt > 0 ) data_cnt = (size_t)cnt;
	} else {
	}
	printf( "Using randomly generated data of size %ld\n", data_cnt );

	KVPair* data = (KVPair*) malloc(sizeof(KVPair)*data_cnt);
	uint64_t mask = 1;
	static const size_t maskbits = 18;
	mask = (mask << maskbits) - 1;
	uint64_t keyhash = 0;
	for ( size_t i = 0; i < data_cnt; i++ ) {
		KVPair p;
		p.key = (rand() & mask);
		p.key <<= maskbits;
		p.key |= (rand() & mask);
		p.val = rand();
		data[i] = p;

		keyhash ^= (p.key + p.val);
	}

	printf("Finished data generation\n"); fflush(stdout);

	std::chrono::high_resolution_clock::time_point start;
	std::chrono::high_resolution_clock::time_point now;
	std::chrono::microseconds duration_micro;
	start = std::chrono::high_resolution_clock::now();
	switch (sorterchoice) {
		case 0:
			quick_sort(data, data_cnt);
			break;
		case 1:
			std::sort(data, data+data_cnt, my_comparator);
			break;
		case 2:
			user_sort(data, data_cnt, threads);
			break;
		default: 
			std::sort(data, data+data_cnt, my_comparator);
			break;
	}
	now = std::chrono::high_resolution_clock::now();
	duration_micro = std::chrono::duration_cast<std::chrono::microseconds> (now - start);
	printf( "Elapsed time: %f s\n", 0.000001f*duration_micro.count() );

	printf("Finished sorting. Evaluation results...\n"); fflush(stdout);


	uint64_t nkhash = 0;
	KeyType lastkey = 0;
	size_t wrongcnt = 0;
	for ( size_t i = 0; i < data_cnt; i++ ) {
		KVPair p = data[i];
		nkhash ^= (p.key + p.val);
		if ( lastkey > p.key ) {
			//printf( "Order wrong\n" );
			wrongcnt++;
		}
		lastkey = p.key;
	}

	if ( wrongcnt > 0 ) {
		printf( "WRONG RESULTS: Key order is wrong!\n" );
	} else if ( keyhash != nkhash ) {
		printf ( "WRONG RESULTS: Key order is correct, but key-value match may be changed\n" );
	} else {
		printf( "CORRECT RESULTS!\n" );
	}
}
