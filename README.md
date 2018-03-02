test-atomic-writes
==================
Test your filesystem's ability to correctly serialize writes as expected in APPEND mode on a POSIX system.
## Tools and ideas:
- [Visit Golang!](https://golang.org)
- [Not the Wizard!](https://www.notthewizard.com/2014/06/17/are-files-appends-really-atomic)
- [Stackoverflow file-append atomic question](http://stackoverflow.com/questions/1154446/is-file-append-atomic-in-unix)
## Issue
- Write in append mode in Mac OS X APFS doesn't appear to be atomic
## Dev log
### 2.24.2018
- Published to Github
- Reported [Apple Bug #37859698](https://bugreport.apple.com/web/?problemID=37859698)
### 2.19.2018 Tested
- MacOS APFS Volume • APFS (Encrypted) -- Fails
```
mikes-air:test-atomic-writes mike$ ./test-atomic-writes -n 100
Each line of file out.tmp will be 4096 bytes, written by 50 workers, writing 100 lines each.
validate: filesize missmatch: filesize=20467712 expected file size=20480000
validate: missed messages: msgCnt=4997 expected=5000 difference=3
validate: source changes=60 shuffle=1.20 msgCnt=4997 errCnt=2
mikes-air:test-atomic-writes mike$ df -T apfs
Filesystem   512-blocks      Used Available Capacity iused               ifree %iused  Mounted on
/dev/disk1s1  489825072 220048920 262175048    46% 1468208 9223372036853307599    0%   /
/dev/disk1s4  489825072   6291496 262175048     3%       3 9223372036854775804    0%   /private/var/vm
```
- MAcOS USB External Physical Volume • Mac OS Extended (Journaled) -- Works
- AWS FreeBSD Local Volume • ufs -- Works
- AWS Amazon-Linux Local Volume • ext4 -- Works
```
Amazon_Linux-512MB-Ohio-1
18.219.2.238
[ec2-user@ip-172-26-2-243 test-atomic-write]$ ./test-atomic-write -n 1000 -s 15013 -w 100
Each line of out.tmp will be 15013 characters long, writen by 100 workers writing 1000 lines each.
validate: fd=3 filesize=1501300000 efs=1501300000
validate: source changes=636 shuffle=6.36 msgCnt=100000 errCnt=0
[ec2-user@ip-172-26-2-243]$ df -T
Filesystem     Type     1K-blocks    Used Available Use% Mounted on
devtmpfs       devtmpfs    241544      56    241488   1% /dev
tmpfs          tmpfs       250508       0    250508   0% /dev/shm
/dev/xvda1     ext4      20509288 3111320  17297720  16% /
[ec2-user@ip-172-26-2-243 ~]$
```
```
FreeBSD-512MB-Ohio-1
18.219.169.234
[ec2-user@ip-172-26-8-133 /usr/home/ec2-user/src/test-atomic-write]$ ./test-atomic-write -n 1000 -s 15913 -w 100
Each line of file out.tmp will be 15913 bytes, writen by 100 workers, writing 1000 lines each.
validate: source changes=86299 shuffle=862.99 msgCnt=100000 errCnt=0
[ec2-user@ip-172-26-8-133 /usr/home/ec2-user/src/test-atomic-write]$
```
### 2.18.2018
- Tested on AWS lightsail linux. Append was unfailing at all sizes and counts up to 16k.
- Tested on AWS FreeBSD lightsail too.
```
mikes-air:test-atomic-write mike$ uname -a
Darwin mikes-air.local 17.4.0 Darwin Kernel Version 17.4.0: Sun Dec 17 09:19:54 PST 2017; root:xnu-4570.41.2~1/RELEASE_X86_64 x86_64
dtruss output.
open("out.tmp\0", 0x20A, 0x180)		 = 3 0
```
- Created FreeBSD linux
```
[ec2-user@ip-172-26-8-133 /usr/home/ec2-user/src/test-atomic-write]$ uname -a
FreeBSD ip-172-26-8-133 11.1-RELEASE-p4 FreeBSD 11.1-RELEASE-p4 #0: Tue Nov 14 06:12:40 UTC 2017     root@amd64-builder.daemonology.net:/usr/obj/usr/src/sys/GENERIC  amd64
[ec2-user@ip-172-26-8-133 /usr/home/ec2-user/src/test-atomic-write]$
[ec2-user@ip-172-26-8-133 /usr/home/ec2-user/src/test-atomic-write]$ ./test-atomic-write -
s 500013 -w 100
Each line of file out.tmp will be 500013 bytes, writen by 100 workers, writing 50 lines each.
validate: source changes=4745 shuffle=47.45 msgCnt=5000 errCnt=0
[ec2-user@ip-172-26-8-133 /usr/home/ec2-user/src/test-atomic-write]$
```
- Created a lightsail instance on AWS to compare test-atomic-write results between MacOS and Linux.
```
Amazon_Linux-512MB-Ohio-1
18.219.2.238
[ec2-user@ip-172-26-2-243 test-atomic-write]$ uname -a
Linux ip-172-26-2-243 4.9.76-3.78.amzn1.x86_64 #1 SMP Fri Jan 12 19:51:35 UTC 2018 x86_64 x86_64 x86_64 GNU/Linux
[ec2-user@ip-172-26-2-243 test-atomic-write]$ ./test-atomic-write -n 10000 -s 16000 &
[1]+  Running                 ./test-atomic-write -n 10000 -s 16000 &
Each line of out.tmp will be 16000 characters long, writen by 50 workers writing 10000 lines each.
[ec2-user@ip-172-26-2-243 test-atomic-write]$ validate: fd=3 filesize=800000000 efs=800000000
[ec2-user@ip-172-26-2-243 test-atomic-write]$ validate: source changes=3431 shuffle=68.62 msgCnt=500000 errCnt=0
[1]+  Done                    ./test-atomic-write -n 10000 -s 16000
```
### 2.13.2018
-  Cleaned up stats. test works, but results seem to suggest that whenever the writes are shuffled between writers, writes are lost.
### 2.10.2018
- writing an append write test program. Modeled after Oz Soloman's shell script for the same purpose. His can be found [here](https://www.notthewizard.com/2014/06/17/are-files-appends-really-atomic/).
