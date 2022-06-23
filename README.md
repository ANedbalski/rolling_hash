# rolling_hash

Application made as a part of test task

## usage
to calculate signature sum:
>rdiff (signature|s) old-file signature-file 

to produce delta file:
> rdiff (delta|d) signature-file new-file delta-file

For the task it uses Adler32 to calculate rolling hash and md5 for the strong hash
In order to extend or update hashing algorythms it made in the separate package. So it can easily be updated to md4 to increase speed or blake32 to decrease collisions

Functionality to store the data is also moved to a separate package. In purpose to change the file format
!For the testing purposes the delta file is storing in the human-readable format. it can be replaced by the binary format make the file smaller
