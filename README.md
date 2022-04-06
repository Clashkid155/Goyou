## Goyou


Simple telegram bot for downloading YouTube media written in GO.


**NOTE:**
- File upload is limited to 20mb 🥲. I 'll fix this soon(hopefully).
- Put bot token in your env or .env file, `botToken = ######.....`.
- Due to rate limiting by YouTube, speed per file is drastically reduced. Issue at kkdai/youtube#232.

Addressing the 20mb limit and slow speed, I'll change the current implementation, so we can have a 2gb limit or use an external server.
For the speed limit, in the issue linked above, 89z indicated a fix which i plan on using, or I can wait for the current library to get updated.

### TODO:
* Playlist support.
* Remove duplicate - filter video formats.
* Speed limit fix.
* Bump upload limit.
