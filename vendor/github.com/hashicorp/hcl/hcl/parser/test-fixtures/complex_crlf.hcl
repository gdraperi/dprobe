variable "foo" ***REMOVED***
	default = "bar"
	description = "bar"
***REMOVED***

variable "groups" ***REMOVED*** ***REMOVED***

provider "aws" ***REMOVED***
	access_key = "foo"
	secret_key = "bar"
***REMOVED***

provider "do" ***REMOVED***
	api_key = "$***REMOVED***var.foo***REMOVED***"
***REMOVED***

resource "aws_security_group" "firewall" ***REMOVED***
	count = 5
***REMOVED***

resource aws_instance "web" ***REMOVED***
	ami = "$***REMOVED***var.foo***REMOVED***"
	security_groups = [
		"foo",
		"$***REMOVED***aws_security_group.firewall.foo***REMOVED***",
		"$***REMOVED***element(split(\",\", var.groups)***REMOVED***",
	]
	network_interface = ***REMOVED***
		device_index = 0
		description = "Main network interface"
	***REMOVED***
***REMOVED***

resource "aws_instance" "db" ***REMOVED***
	security_groups = "$***REMOVED***aws_security_group.firewall.*.id***REMOVED***"
	VPC = "foo"
	depends_on = ["aws_instance.web"]
***REMOVED***

output "web_ip" ***REMOVED***
	value = "$***REMOVED***aws_instance.web.private_ip***REMOVED***"
***REMOVED***
