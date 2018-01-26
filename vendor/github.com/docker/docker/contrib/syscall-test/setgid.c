#include <sys/types.h>
#include <unistd.h>
#include <stdio.h>

int main() ***REMOVED***
	if (setgid(1) == -1) ***REMOVED***
		perror("setgid");
		return 1;
	***REMOVED***
	return 0;
***REMOVED***
