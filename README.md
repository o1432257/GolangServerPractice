# GolangServerPractice
 
## 第一部分:建構基礎Server

### 1.創建一個server類
包含 ip 跟 端口 兩個屬性
### 2.server類創建三個方法
創建一個server對象 NewServer()

啟動server服務 Start()

處理連結業務 Handler()


## 第二部分:用戶上線及廣播功能
### 1.創建一個user類

### 2.user類新增兩個方法
創建一個對象

監聽user對應的channel消息

### 3.server類 新增屬性 
在線用戶的列表   OnlineMap

消息廣播的channel   Message
### 4.server類 在處理客戶端上線的 Handler上創建並添加用戶


### 5.server類 新增兩個方法
廣播消息 BoardCast()

監聽消息 ListenMessage()
### 6.用一個goroutine單獨監聽Message