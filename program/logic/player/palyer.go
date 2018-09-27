package player

import (
	"github.com/gorilla/websocket"
	"sync"
	"chessSever/program/logic/game/poker"
	"strconv"
	"github.com/tidwall/gjson"
	"chessSever/program/logic/msg"
)

/**
定义游戏玩家对象
*/
type Player struct {
	Id       int             //注册用户id
	NickName string          //用户昵称
	Conn     *websocket.Conn //用户socket链接
	HeadPic  string          //用户头像
	Table    *Table          //桌子
	sync.RWMutex
	PokerCards []*poker.PokerCard //玩家手里的扑克牌
	Index    int             //在桌子上的索引
	IsReady  bool			 //是否准备
	IsAuto   bool            //是否托管
}

func NewPlayer(id int, nickName string, conn *websocket.Conn, headPic string) *Player {
	player := Player{
		Id:       id,
		NickName: nickName,
		HeadPic:  headPic,
		Conn:     conn,
	}
	return &player
}

//按照桌号加入牌桌
func (p *Player) JoinTable(key string) error {
	err := getRoom().getTable(key).addPlayer(p)
	return err
}

//开牌桌
func (p *Player) CreateTable(gameName string) {
	table := newTable(p, gameName)
	p.Lock()
	p.Table = table
	p.Unlock()
}

func (p *Player) LeaveTable() {
	p.Table.removePlayer(p)
	p.Lock()
	p.Table = nil
	p.Unlock()
}
//用户跟该桌所有人说话
func (p *Player) SayToTable(msg []byte){
	p.Table.RLock()
	for _,po := range p.Table.Players{
		if po != p {
			po.Conn.WriteMessage(websocket.TextMessage,msg)
		}
	}
	p.Table.RUnlock()
}
//用户跟该桌某一个说话
func (p *Player) sayToAnother(id int,msgB []byte){
	p.Table.RLock()
	for _,po := range p.Table.Players{
		if po.Id == id {
			po.Conn.WriteMessage(websocket.TextMessage,msgB)
		}
	}
	p.Table.RUnlock()
}

func (p *Player)ResolveMsg(msgB []byte) error{
	msgType,err := strconv.Atoi(gjson.Get(string(msgB),"msgType").String())
	if err != nil{
		return err
	}

	switch strconv.Itoa(msgType) {
		case msg.TypeOfAuto:

		case msg.TypeOfUnReady:
			p.unReady()
		case msg.TypeOfReady:
			p.Ready()
		case msg.TypeOfPlay:

		case msg.TypeOfPass:

		case msg.TypeOfLeaveTable:

		case msg.TypeOfJoinTable:

		case msg.TypeOfHint:
	}

	return nil
}

func (p *Player)Ready(){
	p.Lock()
	p.IsReady = true
	p.Unlock()
	p.Table.userReady()
}

func (p *Player)unReady(){
	p.Lock()
	p.IsReady = false
	p.Unlock()
}