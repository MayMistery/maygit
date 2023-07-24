# maygit

1 可以代替缺省值
- mgit -u 可以将时间戳转化为当前的时间
- -t 可以测试ssh和scp连接是否可以成功，通过configure文件中的账号密码，configure文件格式如下
```
host = 127.0.0.1
port = 22
user = root
pass = password
scp = false
```
- `mgit -i` 在当前目录初始化mgit，增加.gitignore，忽略mgit，edr,bk和patch等目录

- `mgit -c “fix ***”`在本地执行git commit，且commit message为“fix ***”，并输出commit_sha

- `mgit -b "/var/www/html,tmp"`可以先连接ssh，然后将远端的/var/www/html目录打包为html_{{时间戳}}.tar.gz到/tmp目录下。var/www/html,/tmp为缺省值。

- `mgit -bk "/var/www/html,tmp"`，可以先连接ssh，然后将远端的/var/www/html目录打包为html_{{时间戳}}.tar.gz到/tmp目录下，然后通过scp或sftp下载到本地的bk目录下。 /var/www/html,/tmp为缺省值。

- `mgit -gen n或{{commit_sha}}` 根据前n次的commit(如果给出commit_sha,则是commit_sha的commit和当前的做对比)和当前的commit的diff信息，制作可以在远端服务器上直接运行的patch脚本，命名为patch_{{时间戳}}.sh到/patch目录下

- `mgit -p "*.sh,/var/www/html"` ,连接远程ssh，scp或sftp上传*.sh并在/var/www/html执行当前patch文件进行服务热重载。",/var/www/html"为缺省值,无*.sh时默认为时间戳最新的.sh文件

- `mgit -emerg "html,/var/www/html,/tmp"` 将html目录打包为.tar.gz，并sftp上传到ssh连接的远端服务器的/tmp目录下，并解压到/var/tmp/html。html,var/www/html,/tmp为缺省值。

- `mgit -hard "*.tar.gz,/var/www/html,/tmp"` 将bk文件夹内的*.tar.gz通过sftp传到ssh连接的远端服务器的/tmp目录下，并解压到/var/www/html。",/var/www/html"为缺省值,无*.tar.gz时默认为时间戳最新的文件

- `mgit -edr 192.168.1.1` 将/edr目录中的全部内容连接ssh并scp上传到远端/tmp目录下，并给edr执行权限

- `mgit -die "/var/www/html"` 连接ssh并强制删除远端所有/var/www/html目录中的内容。/var/www/html为缺省值。

- `mgit -h` 输出帮助信息
