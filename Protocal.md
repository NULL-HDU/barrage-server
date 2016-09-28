# Barrage Protocal

*significant '<f>' means be used for future or it may not be used.*

## struct
*the value in parens is the sturct of this unit*

**userId**: `userId(Uint64)`

**bool**: `bool(Uint8)`

* bool: Uint8, 0 (0x00): false, 1 (0x01): true. 

**damage**: `damage(Uint16)`

**location**: `x(float32) + y(float32)`

* x: float32, x axis.
* y: float32, y axis.

**ballId**: `userId(userId)+id(Uint16)`

* userId: Uint64, the id of the creator of this ball
* id: Uint16, it is a value from 1 - 2*32. After user creating a ball, it add to one. 0 is your airplane

**ball**: `camp(userId) + ballId(ballId) + ballType(Uint8) + hp(Uint16) + damage(damage) + role(Uint8) + special(Uint16) + speed(Uint8) + attackDir(Float32) + alive(bool) + isKilled(bool) + locationCurrent(location)`

**nickname**: `lengthOfName(Uint8) + name(lengthOfNickname * Uint8)`

* lengthOfName: Uint8, the length of nickname
* name: lengthOfNickname * byte, it is a string.

**background**: `imageId(Uint8)`

**collisionInfo**: `ballA(ballId) + ballB(ballId) + damageToA(damage) + damageToB(damage)`

**collisionInfos**: `lengthOfCollisionInfos(Uint32) + collisionInfoArray(lengthOfCollisionInfos * collisionInfo)`

* lengthOfCollisionInfos: Uint32, the length of collisionInfo.

**displacementInfo**: `displacementOfBall(ball)`

**displacementInfos**: `lengthOfDisplacementInfos(Uint32) + displacementInfoArray(lengthOfDisplacementInfos * displacementInfo)`

* lengthOfDisplacementInfos: Uint32, the length of displacementInfo.

**newBallsInfo**: `newBall(ball)`

**newBallsInfos**: `lengthOfNewBallsInfos(Uint32) + newBallsInfoArray(lengthOfNewBallsInfos * newBallsInfo)`

## model 

`message length(Uint32) + timestamp(Int64) + message type(Uint8) + message body(this is different struct from different message types)`

* message length: Uint32, the length of message, including 'length', 'type' and 'body', the unit of length is bytes.
* timestamp: a Unix time, the number of seconds elapsed since January 1, 1970 UTC.
* message type: Uint8, the type of message, type defines the way to decoding the message and what should ends do.
* message body: this is different struct from message which has different type.

## Client send to Server

### <f>1. enter room

type value: 1  (0x01)

message body: `userId(Uint64) + nickname(nickname) + roomNumber(Uint32) + troop(Uint8)`

* userId: Uint64, the id of uint64.
* nickname: nickname, the name of user
* roomNumber: Uint32, the room of game.
* troop: Uint8, the troop of user.

### <f>2. ready game

type value: 2  (0x02)

message body: `userId(Uint64) + roomNumber(Uint32) + readyValue(Uint8)`

* userId: Uint64, the id of user.
* roomNumber: Uint32, the room of game.
* readyValue: Uint8, 0 (0x00): cancel, 1 (0x01): ready.

### <f>3. start gamme

type value: 3  (0x03)

message body: `userId(Uint64) + roomNumber(Uint32)`

* userId: Uint64, the id of uint64.
* roomNumber: Uint32, the room of game.

### 8. disconnect(leave early)

type value: 8  (0x08)

message body: `userId(Uint64) + roomNumber(Uint32)`

* userId: Uint64, the id of uint64.
* roomNumber: Uint32, the room of game.

### 9. connect(join game)

type value: 9  (0x09)

message body: `userId(Uint64) + nickname(nickname) + roomNumber(Uint32) + troop(Uint8)`

* userId: Uint64, the id of uint64.
* nickname: nickname, the name of user
* roomNumber: Uint32, the room of game.
* troop: Uint8, the troop of user.


## Server send to Client

### <f>4. someone ready

type value: 4  (0x04)

message body: `readyUserId(Uint64) + roomNumber(Uint32) + readyValue(Uint8)`

* readyUserId: Uint64, the id of the ready one.
* roomNumber: Uint32, the room of game.
* readyValue: Uint8, 0 (0x00): cancel, 1 (0x01): ready.

### <f>5. game starts

type value: 5  (0x05)

message body: `roomNumber(Uint32)`

* roomNumber: Uint32, the room of game.

### 6. airplane of userself has been created

type value: 6  (0x06)

message body: `airplane(ball) + nickname(nickname)`

* airplane: ball, user's airplane
* nickname: nickname, the name of user

### 7. playground info

type value: 7  (0x06)

message body: `collisionInfos(collisionInfos) + displacementInfos(displacementInfos) + newBallsInfos(newBallsInfos)`

* collisionInfos: collisionInfos, the information about ball collision.
* displacementInfos: displacementInfos, the information about balls displacement, include death info and disappear info.
* newBallsInfos: newBallsInfos, the information about new balls.

### 10. special message

type value: 10  (0x0a)

message body: `lengthOfSpecialMessage(Uint8) + specialMessage(lengthOfSpecialMessage * Uint8)`