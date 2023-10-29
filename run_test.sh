#!/bin/bash
#

host -T poppy.home 127.0.0.1
host -T 10.0.0.8 127.0.0.1

host -T a_record_1.home 127.0.0.1
host -T 192.168.1.1 127.0.0.1

host poppy.home 127.0.0.1
host 10.0.0.8 127.0.0.1

host a_record_1.home 127.0.0.1
host 192.168.1.1 127.0.0.1
