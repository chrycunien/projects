#include <stdio.h>
#include <tuple>
#include "types.h"
// for pthread affinity
#include <sched.h>
#include <unistd.h>

#include <pthread.h>
#include <iostream>
#include <vector>
#include <stack>
#include <queue>
#include <limits.h>
using namespace std;

bool my_comparator( KVPair a, KVPair b)
{	
		return a.key < b.key;
}

size_t partition(KVPair* array, size_t data_cnt) {
	KeyType p = array[data_cnt-1].key;

	size_t curidx = 0;
	for ( size_t i = 0; i < data_cnt-1; i++ ) {
		if ( array[i].key <= p ) {
			KVPair t = array[i];
			array[i] = array[curidx];
			array[curidx] = t;
			curidx++;
		}
	}

	KVPair t = array[curidx];
	array[curidx] = array[data_cnt-1];
	array[data_cnt-1] = t;


	return curidx;
}
void quick_sort(KVPair* array, size_t data_cnt) {
	size_t pivot = partition(array, data_cnt);
	if ( pivot > 1 ) quick_sort(array, pivot);
	if ( pivot < data_cnt - 1 ) quick_sort(array+pivot, data_cnt-pivot);
}


// int stick_this_thread_to_core(int core_id) {
//    int num_cores = sysconf(_SC_NPROCESSORS_ONLN);
//    if (core_id < 0 || core_id >= num_cores)
//       return EINVAL;

//    cpu_set_t cpuset;
//    CPU_ZERO(&cpuset);
//    CPU_SET(core_id, &cpuset);

//    pthread_t current_thread = pthread_self();    
//    return pthread_setaffinity_np(current_thread, sizeof(cpu_set_t), &cpuset);
// }

typedef struct sort_arg {
	KVPair *arr;
	size_t size;
	int id;
} sort_arg_t;

typedef struct merge_sort_arg {
	KVPair *arr;
	size_t size;
	size_t cnt;
} merge_sort_arg_t;

typedef pair<KVPair*, size_t> pair_t;
// typedef tuple<size_t, size_t, size_t> range_t;

class SortedPartition {
	public:
		SortedPartition(KVPair* arr, size_t size) {
			this->arr = arr;
			this->size = size;
		}

		bool is_empty() {
			return this->size == 0;
		}

		KVPair top() {
			return this->arr[0];
		}

		KVPair next() {
			KVPair elem = this->arr[0];
			this->arr++;
			this->size--;
			return elem;
		}

		void print() {
			printf("arr: %p, size: %zu\n", arr, size);
		}
	
	private:
		KVPair* arr;
		size_t size;
};

class KVPairComparator 
{ 
public: 
	int operator() (const KVPair& p1, const KVPair& p2) 
	{ 
		return p1.key > p2.key; 
	} 
}; 

class SortedMerger {
	public:
		SortedMerger(KVPair* arr, size_t size, int k) {
			int par_size = size / k;

			for (int i = 0; i < k; i++) {
				if (i == k - 1) {
					list.push_back(SortedPartition(arr + i * par_size, size - par_size * i));
				} else {
					list.push_back(SortedPartition(arr + i * par_size, par_size));
				}
			}

			// for (int i = 0; i < list.size(); i++) {
			// 	list[i].print();
			// }

			// assume all elements are good
			for (int i = 0; i < list.size(); i++) {
				q.push(list[i].next());
			}

			this->clean_flag = false;

			// for (int i = 0; i < list.size(); i++) {
			// 	list[i].print();
			// }

			printf("start merging...\n");
		}


		bool is_empty() {
			return q.size() == 0;
		}

		KVPair next() {
			KVPair elem = q.top(); q.pop();
			// printf("** key: %llu, value: %u\n", elem.key, elem.val);
			if (list.size() > 0) {
				fetch_next();
			}
			return elem;
		}

	private:

		vector<SortedPartition> list;
		// min heap
		priority_queue<KVPair, vector<KVPair>, KVPairComparator> q;
		bool clean_flag;
		int clean_index;

		
		void fetch_next() {
			size_t min_value = INT64_MAX;
			size_t min_index = -1;
			for (int i = 0; i < list.size(); i++) {
				// printf("+ i: %d\n", i);
				if (list[i].is_empty()) {
					this->clean_flag = true;
					this->clean_index = i;
					continue;
				}
				if (list[i].top().key < min_value) {
					min_index = i;
					min_value = list[i].top().key;
				}
			}
			if (list.size() > 1 || !clean_flag) {
				// printf("min_index: %lu\n", min_index);
				KVPair elem = list[min_index].next();
				// printf("++ key: %llu, value: %u\n", elem.key, elem.val);
				q.push(elem);
			}

			if (clean_flag) {
				this->clean();
				clean_flag = false;
			}
		}

		void clean() {
			list.erase(list.begin() + clean_index);
		}
};

class LazyMerger {
	public:
		LazyMerger(KVPair* arr, size_t size, int k) {
			this->arr_size = size;
			this->sorted_idx = 0;

			if (size <= 1000) {
				quick_sort(arr, size);
				this->sorted_arr = arr;
				this->is_leaf = true;
				return ;
			}

			this->is_leaf = false;

			int par_size = size / k;

			for (int i = 0; i < k; i++) {
				if (i == k - 1) {
					list.push_back(LazyMerger(arr + i * par_size, size - par_size * i, k));
				} else {
					list.push_back(LazyMerger(arr + i * par_size, par_size, k));
				}
			}

			// for (int i = 0; i < list.size(); i++) {
			// 	list[i].print();
			// }

			// assume all elements are good
			for (int i = 0; i < list.size(); i++) {
				q.push(list[i].next());
			}

			this->clean_flag = false;

			// for (int i = 0; i < list.size(); i++) {
			// 	list[i].print();
			// }

			// printf("start merging...\n");
		}

		bool is_empty() {
			if (is_leaf) {
				return sorted_idx == arr_size;
			} else {
				return q.size() == 0;
			}
		}

		KVPair peek() {
			if (is_leaf) {
				return sorted_arr[sorted_idx];
			} else {
				return q.top();
			}
		}

		KVPair next() {
			if (is_leaf) {
				return sorted_arr[sorted_idx++];
			} else {
				KVPair elem = q.top(); q.pop();
				// printf("** key: %llu, value: %u\n", elem.key, elem.val);
				if (list.size() > 0) {
					fetch_next();
				}
				return elem;
			}
		}

	private:
		KVPair* sorted_arr;
		int sorted_idx;
		int arr_size;
		bool is_leaf;
	
		vector<LazyMerger> list;
		// min heap
		priority_queue<KVPair, vector<KVPair>, KVPairComparator> q;
		bool clean_flag;
		int clean_index;

		
		void fetch_next() {
			size_t min_value = INT64_MAX;
			size_t min_index = -1;
			for (int i = 0; i < list.size(); i++) {
				// printf("+ i: %d\n", i);
				if (list[i].is_empty()) {
					this->clean_flag = true;
					this->clean_index = i;
					continue;
				}
				if (list[i].peek().key < min_value) {
					min_index = i;
					min_value = list[i].peek().key;
				}
			}
			if (list.size() > 1 || !clean_flag) {
				// printf("min_index: %lu\n", min_index);
				KVPair elem = list[min_index].next();
				// printf("++ key: %llu, value: %u\n", elem.key, elem.val);
				q.push(elem);
			}

			if (clean_flag) {
				this->clean();
				clean_flag = false;
			}
		}

		void clean() {
			list.erase(list.begin() + clean_index);
		}
};

void *inner_quick_sort(void *argi) {
	sort_arg_t *arg = (sort_arg_t *) argi;
	quick_sort(arg->arr, arg->size);
	pthread_exit(NULL);
}


pair_t inner_merge_sort2(pair_t p1, pair_t p2) {
	size_t l1 = p1.second, l2 = p2.second, k = 0, i1 = 0, i2 = 0;
	KVPair *arr1 = p1.first, *arr2 = p2.first;
	KVPair* new_arr = new KVPair[l1 + l2];

	while (i1 < l1 && i2 < l2) {
		KVPair e1 = arr1[i1], e2 = arr2[i2];
		if (e1.key < e2.key) {
			new_arr[k] = e1;
			i1++; k++;
		} else {
			new_arr[k] = e2;
			i2++; k++;
		}
	}

	// append the remaining items
	while (i1 < l1) {
		new_arr[k] = arr1[i1];
		i1++; k++;
	}
	while (i2 < l2) {
		new_arr[k] = arr2[i2];
		i2++; k++;
	}

	return make_pair(new_arr, l1 + l2);
}

KVPair* inner_k_merge_sort(KVPair* arr, size_t size, int k) {
	KVPair* new_arr = new KVPair[size];
	SortedMerger merger(arr, size, k);
	int i = 0;

	while (!merger.is_empty()) {
		KVPair elem = merger.next();
		// printf("(%d) key: %llu, value: %u\n", i, elem.key, elem.val);
		new_arr[i++] = elem; 
	}

	return new_arr;
}

void k_merge_sort(KVPair* arr, size_t size, int k) {
	KVPair* new_arr = inner_k_merge_sort(arr, size, k);

	// direct copy
	memcpy(arr, new_arr, sizeof(KVPair) * size);
}

KVPair *inner_merge_sort_2(KVPair* arr, size_t size, int threads_num) {
	KVPair* new_arr = new KVPair[size];
	int par_size = size / threads_num;
	queue<pair_t> q;
	for (int i = 0; i < threads_num; i++) {
		if (i == threads_num - 1) {
			q.push(make_pair(arr + i * par_size, size - par_size * i));
		} else {
			q.push(make_pair(arr + i * par_size, par_size));
		}
		// cout << q.back().first << " " << q.back().second << endl;
	}

	while (q.size() > 1) {
		pair_t p1 = q.front(); q.pop();
		pair_t p2 = q.front(); q.pop();
		// printf("** arr1: %p, size1: %zu. arr2: %p, size2: %zu.\n", p1.first, p1.second, p2.first, p2.second);
		pair_t merged_pair = inner_merge_sort2(p1, p2);
		q.push(merged_pair);
	}

	new_arr = q.front().first;

	return new_arr;
}

void merge_sort_2(KVPair* arr, size_t size, int k) {
	KVPair* new_arr = inner_merge_sort_2(arr, size, k);

	// direct copy
	memcpy(arr, new_arr, sizeof(KVPair) * size);
}

vector<size_t> get_count(int k, int slots) {
	vector<size_t> count(slots);
	for (int i = 0; i < k; i++) {
		count[i % slots]++;
	}
	return count;
}

// queue<range_t> init_queue(vector<size_t> count, size_t size, size_t par_size) {
// 	queue<range_t> q;
// 	size_t acc = 0;
// 	size_t k = count.size();
// 	for (int i = 0; i < count.size(); i++) {
// 		if (i == count.size() - 1) {
// 			q.push(make_tuple(acc, size - acc, count[i]));
// 		} else {
// 			q.push(make_tuple(acc,  par_size * count[i], count[i]));
// 			acc += par_size * count[i];
// 		}
// 	}
// 	return q;
// }


// KVPair *merge_sort_k_merger(KVPair* arr, size_t size, int threads_num) {
// 	KVPair* new_arr = new KVPair[size];
// 	int par_size = size / threads_num;
// 	vector<size_t> count = get_count(threads_num, threads_num / 2);
// 	queue<range_t> q = init_queue(count, size, par_size);

// combine to 2^n
// then do static dispatching

// 	while (q.size() > 1) {
// 		range_t e = q.front(); q.pop();
// 		size_t start = get<0>(e), len = get<1>(e), par = get<2>(e);
		
// 		KVPair* merged_pair = inner_k_merge_sort(arr + start, len, par);
// 		q.push(make_tuple(start, ));
// 	}

// 	new_arr = q.front().first;

// 	return new_arr;
// }



KVPair* inner_funnel_sort(KVPair* arr, size_t size) {
	if (size <= 1000) {
		quick_sort(arr, size);
		return arr;
	}

	KVPair* new_arr = new KVPair[size];
	LazyMerger merger(arr, size, 10);
	int i = 0;

	while (!merger.is_empty()) {
		KVPair elem = merger.next();
		// printf("(%d) key: %llu, value: %u\n", i, elem.key, elem.val);
		new_arr[i++] = elem; 
	}

	return new_arr;
}

void *funnel_sort(void *argi) {
	sort_arg_t* arg = (sort_arg_t*) argi;
	KVPair* arr = arg->arr;
	size_t size = arg->size;
	if (size <= 32000) {
		quick_sort(arr, size);
	} else {
		KVPair* new_arr = inner_funnel_sort(arr, size);
		memcpy(arr, new_arr, sizeof(KVPair) * size);
	}
	pthread_exit(NULL);
}


void user_sort(KVPair* array, size_t data_cnt, size_t threads_num) {
	if (data_cnt < 64000) {
		quick_sort(array, data_cnt);
		return ;
	}

	int num_cores = sysconf(_SC_NPROCESSORS_ONLN);

	threads_num = num_cores <= threads_num ? num_cores : threads_num;

	pthread_t threads[threads_num];
	sort_arg_t td[threads_num];
	int par = data_cnt / threads_num;
	pthread_attr_t attr;
	void *status;

	// Initialize and set thread joinable
	pthread_attr_init(&attr);
	pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_JOINABLE);

	for(int i = 0; i < threads_num; i++) {
		td[i].arr = array + par * i;
		td[i].size = (i == threads_num - 1) ? data_cnt - (threads_num - 1) * par : par;
		td[i].id = i;

		if (pthread_create(&threads[i], &attr, inner_quick_sort, (void *)&td[i]) != 0) {
			cout << "Error: unable to create thread." << endl;
			exit(-1);
		}

		// if (pthread_create(&threads[i], &attr, funnel_sort, (void *)&td[i]) != 0) {
		// 	cout << "Error: unable to create thread." << endl;
		// 	exit(-1);
		// }
   }

	// Free attribute and wait for the other threads
	pthread_attr_destroy(&attr);

	for(int i = 0; i < threads_num; i++) {
		if (pthread_join(threads[i], &status) != 0) {
			cout << "Error: unable to join." << endl;
			exit(-1);
		}
	}

	k_merge_sort(array, data_cnt, threads_num);
	// merge_sort_2(array, data_cnt, threads_num);

	// funnel_sort(array, data_cnt);

}