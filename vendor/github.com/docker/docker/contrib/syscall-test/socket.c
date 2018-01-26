#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

int main() ***REMOVED***
	int s;
	struct sockaddr_in sin;

	s = socket(AF_INET, SOCK_STREAM, 0);
	if (s == -1) ***REMOVED***
		perror("socket");
		return 1;
	***REMOVED***

	sin.sin_family = AF_INET;
	sin.sin_addr.s_addr = INADDR_ANY;
	sin.sin_port = htons(80);

	if (bind(s, (struct sockaddr *)&sin, sizeof(sin)) == -1) ***REMOVED***
		perror("bind");
		return 1;
	***REMOVED***

	close(s);

	return 0;
***REMOVED***
