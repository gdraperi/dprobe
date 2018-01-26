# ESCAPE=\



FROM brimstone/ubuntu:14.04


# TORUN -v /var/run/docker.sock:/var/run/docker.sock


ENV GOPATH \
/go



# Install the packages we need, clean up after them and us
RUN apt-get update \
	&& dpkg -l | awk '/^ii/ ***REMOVED***print $2***REMOVED***' > /tmp/dpkg.clean \


    && apt-get install -y --no-install-recommends git golang ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists \

	&& go get -v github.com/brimstone/consuldock \
    && mv $GOPATH/bin/consuldock /usr/local/bin/consuldock \

	&& dpkg -l | awk '/^ii/ ***REMOVED***print $2***REMOVED***' > /tmp/dpkg.dirty \
	&& apt-get remove --purge -y $(diff /tmp/dpkg.clean /tmp/dpkg.dirty | awk '/^>/ ***REMOVED***print $2***REMOVED***') \
	&& rm /tmp/dpkg.* \
	&& rm -rf $GOPATH




