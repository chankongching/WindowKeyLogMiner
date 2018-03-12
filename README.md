# Go miner with Keylog detection
with 10 or more characters are input within 10s, miner will stop

### Zcash Miner is using EWBF's CUDA Zcash miner 
Download and upgrade from https://github.com/nanopool/ewbf-miner/releases
then unzip "Zcash Miner x.x.x.zip" and place the miner into folder "zcashminer"

## KeyLogger
Key Logger function is referencing
https://github.com/SaturnsVoid/Windows-KeyLogger

### Configuration file input
https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go

### Execution
Running mining.vbs to run in background
Running WindowKeyLogMiner.exe to run the exe in front
Miner will always be ran in background

### Functionality
After running the .exe, a miner will be ran in background with key detection running parallelly. Once 11 continuous key is being typed within 10s, miner will be off
