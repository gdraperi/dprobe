#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/ip.h>
#include <netinet/udp.h>

int main() ***REMOVED***
	if (socket(PF_INET, SOCK_RAW, IPPROTO_UDP) == -1) ***REMOVED***
		perror("socket");
		return 1;
	***REMOVED***

	return 0;
***REMOVED***
