# 说明
一个学习项目，使用golang编写一个简单的扫描器

# 目标
- [x]  利用socket实现全连接的端口扫描功能
- [ ]  实现存活判断
    - [X]  PING调用
    - [X]  ICMP协议 (通用)
    - [ ]  ARP协议  (内网更适用)
- [ ]  不依赖libcap实现syn端口扫描(windows api仅适用win平台)
- [ ]  实现弱口令扫描
- [ ]  实现端口服务识别
- [ ]  实现web Title获取
- [ ]  实现web 指纹识别
- [ ]  ...

# Update
- 20220105 初始化该学习项目
- 20220105 2:57 实现系统PING调用检测存活IP
- 20220105 12:27 实现并发ping探测存活
- 20220105 16:57 实现ICMP协议探测存活

# Bug
- 20220106 在使用ICMP协议探测B段时会存在误报情况（不存活也报存活）