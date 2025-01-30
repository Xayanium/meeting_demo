### 数据库设计
**用户表：** 存放用户相关信息（用户名、密码、webRTC需要的sdp）

**会议房间表：** 存放房间相关信息（房间uuid号、房间名字、创建时间、结束时间、创建房间用户的id）

**用户和房间关联表：** 将进入房间的用户和房间关联


### Pion WebRTC库用法
1. 创建新API实例

2. 创建新PeerConnection

3. 添加ICE Candidate

4. 创建并设置 SDP offer/answer 为 local description

5. 接受 SDP answer/offer 为 remote description

6. 处理连接成功事件

7. 添加媒体轨道


### 完整的WebRTC实时语音或视频对话的基本流程
1. **初始化RTCPeerConnection**：双方都需要创建RTCPeerConnection实例，并配置ICEServers等信息。

2. **获取本地媒体流**：使用getUserMedia获取音频、视频流，并添加到RTCPeerConnection实例中。

3. **创建和交换SDP Offer/Answer**：一方创建offer，另一方根据offer创建answer，两者之间通过信令服务器交换这些信息。

4. **ICE候选交换**：两端交换ICE候选者信息以穿透NAT和防火墙，确保直接的点对点连接。

5. **媒体流传输**：当连接建立后，媒体流将自动在两个端点之间传输。

6. **处理断开连接等事件**：监听相关事件如iceconnectionstatechange来处理连接状态的变化。


