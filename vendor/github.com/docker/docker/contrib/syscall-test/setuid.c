#include <sys/types.h>
#include <unistd.h>
#include <stdio.h>

int main() ***REMOVED***
	if (setuid(1) == -1) ***REMOVED***
		perror("setuid");
		return 1;
	***REMOVED***
	return 0;
***REMOVED***
