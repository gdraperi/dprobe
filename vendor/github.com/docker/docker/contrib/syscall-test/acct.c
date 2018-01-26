#define _GNU_SOURCE
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

int main(int argc, char **argv)
***REMOVED***
	int err = acct("/tmp/t");
	if (err == -1) ***REMOVED***
		fprintf(stderr, "acct failed: %s\n", strerror(errno));
		exit(EXIT_FAILURE);
	***REMOVED***
	exit(EXIT_SUCCESS);
***REMOVED***
