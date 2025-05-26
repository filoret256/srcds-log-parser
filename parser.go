package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Структура для представления игрока
type Player struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	SteamID string `json:"steamid"`
	Team    string `json:"team"`
	Position []int `json:"position,omitempty"` // Omit if empty
}

// Базовая структура для всех событий
type BaseEvent struct {
	Event string    `json:"event"`
	Date  EventDate `json:"date"`
}

// Структура для временной метки события
type EventDate struct {
	Year    int    `json:"year"`
	Month   int    `json:"month"` // Note: Go's Month is a type, maybe use int for consistency with JS code
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minutes int    `json:"minutes"`
	Seconds int    `json:"seconds"`
	TS      int64  `json:"ts"`      // Unix timestamp
	FullDate time.Time `json:"fulldate"`
}

// Структуры для конкретных типов событий
type TeamNameEvent struct {
	BaseEvent
	Team     string `json:"team"`
	TeamName string `json:"teamName"`
}

type AttackedEvent struct {
	BaseEvent
	PlayerA *Player `json:"playerA"`
	PlayerB *Player `json:"playerB"`
	Weapon  string  `json:"weapon"`
	Damage  int     `json:"damage"`
	DamageArmor int `json:"damage_armor"`
	Health  int     `json:"health"`
	Armor   int     `json:"armor"`
	Hitgroup string `json:"hitgroup"`
}

type KilledEvent struct {
	BaseEvent
	PlayerA   *Player `json:"playerA"`
	PlayerB   *Player `json:"playerB"`
	Weapon    string  `json:"weapon"`
	Headshot  bool    `json:"headshot"`
	Penetrated bool   `json:"penetrated"`
}

type TriggerEvent struct {
	BaseEvent
	TriggerType string `json:"triggerType"`
	Trigger     string `json:"trigger"`
	Team        string `json:"team,omitempty"` // For team trigger
	ResultScoreCT int `json:"resultScoreCT,omitempty"` // For team trigger
	ResultScoreT int `json:"resultScoreT,omitempty"` // For team trigger
	TriggerOn   string `json:"triggerOn,omitempty"` // For world on trigger
	Player      *Player `json:"player,omitempty"` // For player trigger
	Value       string `json:"value,omitempty"` // For player trigger with value
	PlayerA     *Player `json:"playerA,omitempty"` // For kill location trigger
	PlayerB     *Player `json:"playerB,omitempty"` // For kill location trigger
	// Weaponstats fields
	Weapon string `json:"weapon,omitempty"`
	Shots string `json:"shots,omitempty"`
	Hits string `json:"hits,omitempty"`
	Kills string `json:"kills,omitempty"`
	Headshots string `json:"headshots,omitempty"`
	Tks string `json:"tks,omitempty"`
	Damage string `json:"damage,omitempty"`
	Deaths string `json:"deaths,omitempty"`
	// Weaponstats2 fields (using string as in JS, might need conversion)
	Head string `json:"head,omitempty"`
	Chest string `json:"chest,omitempty"`
	Stomach string `json:"stomach,omitempty"`
	Leftarm string `json:"leftarm,omitempty"`
	Rightarm string `json:"rightarm,omitempty"`
	Leftleg string `json:"leftleg,omitempty"`
	Rigthleg string `json:"rigthleg,omitempty"`

}

type PurchaseEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Item   string  `json:"item"`
}

type DisconnectedEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Reason string  `json:"reason"`
}

type EnteredEvent struct {
	BaseEvent
	Player *Player `json:"player"`
}

type LogStartEvent struct {
	BaseEvent
}

type LogCloseEvent struct {
	BaseEvent
}

type ConnectedEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Address string `json:"adress"` // Typo adress from original code preserved
}

type RconEvent struct {
	BaseEvent
	Address string `json:"address"`
	Port    int    `json:"port"`
	Command string `json:"command"`
}

type ThrewEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Item   string  `json:"item"`
}

type PickedUpEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Item   string  `json:"item"`
}

type SuicideEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Item   string  `json:"item"`
}

type ChangedRoleEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Role   string  `json:"role"`
}

type SIDValidatedEvent struct {
	BaseEvent
	Player *Player `json:"player"`
}

type JoinedTeamEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Team   string  `json:"team"`
}

type SwitchTeamEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	TeamA  string  `json:"teamA"`
	TeamB  string  `json:"teamB"`
}

type LoadingMapEvent struct {
	BaseEvent
	Map string `json:"map"`
}

type StartedMapEvent struct {
	BaseEvent
	Map string `json:"map"`
	CRC string `json:"crc"`
}

type ServerCvarEvent struct {
	BaseEvent
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}

type ScoreEvent struct {
	BaseEvent
	Team      string `json:"team"`
	ScoreType string `json:"scoreType"` // "current", "scored", "final"
	Score     int    `json:"score"`
	Players   int    `json:"players"`
}

type SpawnedEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Spawned string `json:"spawned"` // The role/type spawned as
}

type SayEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Text   string  `json:"text"`
}

type SayTeamEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Text   string  `json:"text"`
}

type PositionReportEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Text   string  `json:"text"` // The raw report string
}

type AssistEvent struct {
	BaseEvent
	PlayerA *Player `json:"playerA"`
	PlayerB *Player `json:"playerB"`
}

type VotingEvent struct {
	BaseEvent
	VoteType   string `json:"voteType"` // "kick", "nextlevel", "changelevel", "scrambleteams", "swapteams"
	VoteResult string `json:"voteResult"` // "success", "failed " (preserve space from JS)
	Player     string `json:"player,omitempty"` // For kick votes
	Map        string `json:"map,omitempty"` // For map votes
}

type MolotovProjectileEvent struct {
	BaseEvent
	Position []int `json:"position"`
	PositionVelocity []int `json:"positionVelocity"`
}

// Регулярные выражения (скомпилированы один раз при старте программы)
var (
	reTeamName           = regexp.MustCompile(`^Team playing "(.+)": (.+)$`)
	reAttacked           = regexp.MustCompile(`^"(.+)" \[(.+) (.+) (.+)\] attacked "(.+)" \[(.+) (.+) (.+)\] with "(.+)" \(damage "(.+)"\) \(damage_armor "(.+)"\) \(health "(.+)"\) \(armor "(.+)"\) \(hitgroup "(.+)"\)$`)
	reKilledHs           = regexp.MustCompile(`^"(.+)" \[(.+) (.+) (.+)\] killed "(.+)" \[(.+) (.+) (.+)\] with "(.+)" \(headshot\)$`)
	reKilledPenHs        = regexp.MustCompile(`^"(.+)" \[(.+) (.+) (.+)\] killed "(.+)" \[(.+) (.+) (.+)\] with "(.+)" \(headshot penetrated\)$`)
	reKilledPen          = regexp.MustCompile(`^"(.+)" \[(.+) (.+) (.+)\] killed "(.+)" \[(.+) (.+) (.+)\] with "(.+)" \(penetrated\)$`)
	reKilled             = regexp.MustCompile(`^"(.+)" \[(.+) (.+) (.+)\] killed "(.+)" \[(.+) (.+) (.+)\] with "(.+)"$`)
	reTriggerTeam        = regexp.MustCompile(`^Team "(.+)" triggered "(.+)" \(CT "(.+)"\) \(T "(.+)"\)$`)
	reTriggerWorldOn     = regexp.MustCompile(`^World triggered "(.+)" on "(.+)"$`)
	reTriggerWorld       = regexp.MustCompile(`^World triggered "(\S+)"$`)
	reTriggerPlayer      = regexp.MustCompile(`^"(.+)" triggered "(.+)"$`)
	reTriggerPlayerVal   = regexp.MustCompile(`^"(.+)" triggered "(.+)" \(value "?(.+)"?\)$`)
	rePurchase           = regexp.MustCompile(`^"(.+)" purchased "(.+)"$`)
	reDisconnected       = regexp.MustCompile(`^"(.+)" disconnected \(reason "(.+)"\)$`)
	reEntered            = regexp.MustCompile(`^"(.+)" entered the game$`)
	reStartLog           = regexp.MustCompile(`^Log file started$`)
	reCloseLog           = regexp.MustCompile(`^Log file closed$`)
	reConnected          = regexp.MustCompile(`^"(.+)" connected, address ""$`) // Adjusted regex for empty address
	reConnectedAddress   = regexp.MustCompile(`^"(.+)" connected, address "(.+)"$`)
	reRcon               = regexp.MustCompile(`^rcon from "(.+):(\d+)"[:] command "(.+)"$`)
	reThrew              = regexp.MustCompile(`^"(.+)" threw (.+) \[(.+) (.+) (.+)\]$`)
	rePickedUp           = regexp.MustCompile(`^"(.+)" picked up item "(.+)"$`)
	reSuicide            = regexp.MustCompile(`^"(.+)" committed suicide with "(.+)"$`)
	reSuicidePosition2   = regexp.MustCompile(`^"(.+)" \[(.+) (.+) (.+)\] committed suicide with "(.+)"$`)
	reSuicidePosition    = regexp.MustCompile(`^"(.+)" committed suicide with "(.+)" \(attacker_position "(.+) (.+) (.+)"\)$`)
	reChangedRole        = regexp.MustCompile(`^"(.+)" changed role to "(.+)"$`)
	reSIDValidated       = regexp.MustCompile(`^"(.+)" STEAM USERID validated$`)
	reJoinedTeam         = regexp.MustCompile(`^"(.+)" joined team "(.+)"$`)
	reSwitchTeam         = regexp.MustCompile(`^"(.+)" switched from team <(.+)> to <(.+)>$`)
	reLoadingMap         = regexp.MustCompile(`^Loading map "(.+)"$`)
	reStartedMap         = regexp.MustCompile(`^Started map "(.+)" \(CRC "(.+)"\)$`)
	reServerCvar         = regexp.MustCompile(`^server_cvar: "(.+)" "(.+)"$`)
	reScoreCurrent       = regexp.MustCompile(`^Team "(.+)" current score "(\d+)" with "(\d+)" players$`)
	reScoreScored        = regexp.MustCompile(`^Team "(.+)" scored "(\d+)" with "(\d+)" players$`)
	reScoreFinal         = regexp.MustCompile(`^Team "(.+)" final score "(\d+)" with "(\d+)" players$`)
	reSpawned            = regexp.MustCompile(`^"(.+)" spawned as "(.+)"$`)
	reSay                = regexp.MustCompile(`^"(.+)" say "(.+)"$`)
	reSayTeam            = regexp.MustCompile(`^"(.+)" say_team "(.+)"$`)
	rePositionReport     = regexp.MustCompile(`^"(.+)" position_report (.+)$`)
	reAssist             = regexp.MustCompile(`^"(.+)" assisted killing "(.+)"$`)
	reVoteSuccessKick    = regexp.MustCompile(`^Vote succeeded "Kick (.+)"$`)
	reVoteFalseKick      = regexp.MustCompile(`^Vote failed "Kick (.+)"$`)
	reVoteSuccessNxLvl   = regexp.MustCompile(`^Vote succeeded "NextLevel (.+)"$`)
	reVoteFalseNxLvl     = regexp.MustCompile(`^Vote failed "NextLevel (.+)"$`)
	reVoteSuccessChLvl   = regexp.MustCompile(`^Vote succeeded "ChangeLevel (.+)"$`)
	reVoteFalseChLvl     = regexp.MustCompile(`^Vote failed "ChangeLevel (.+)"$`)
	reVoteSuccessScramble = regexp.MustCompile(`^Vote succeeded "ScrambleTeams "$`)
	reVoteFalseScramble  = regexp.MustCompile(`^Vote failed "ScrambleTeams "$`)
	reVoteSuccessSwap    = regexp.MustCompile(`^Vote succeeded "SwapTeams "$`)
	reVoteFalseSwap      = regexp.MustCompile(`^Vote failed "SwapTeams "$`)
	reMolotov            = regexp.MustCompile(`^Molotov projectile spawned at (.+) (.+) (.+), velocity (.+) (.+) (.+)$`)

	// hlx2 specific regexes
	reTriggerKillLocation = regexp.MustCompile(`^World triggered "(.+)" \(attacker_position "(.+) (.+) (.+)"\) \(victim_position "(.+) (.+) (.+)"\)$`)
	reTriggerWeaponstats = regexp.MustCompile(`^"(.+)" triggered "(.+)" \(weapon "(.+)"\) \(shots "(.+)"\) \(hits "(.+)"\) \(kills "(.+)"\) \(headshots "(.+)"\) \(tks "(.+)"\) \(damage "(.+)"\) \(deaths "(.+)"\)$`)
	reTriggerWeaponstats2 = regexp.MustCompile(`^"(.+)" triggered "(.+)" \(weapon "(.+)"\) \(head "(.+)"\) \(chest "(.+)"\) \(stomach "(.+)"\) \(leftarm "(.+)"\) \(rightarm "(.+)"\) \(leftleg "(.+)"\) \(rightleg "(.+)"\)$`)
)

// parsePlayer парсит строку игрока в структуру Player
func parsePlayer(playerStr string) *Player {
	var result []string

	// Original JS patterns translated
	if result = regexp.MustCompile(`^(.+)<(\d+)><([STEAM1234567890_:]+)><(.+)>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: result[4]}
	}
	if result = regexp.MustCompile(`^(.+)<(\d+)><([STEAM1234567890_:]+)><>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: "Unassigned"}
	}
	if result = regexp.MustCompile(`^(.+)<(\d+)><([STEAM1234567890_:]+)>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: "Unassigned"}
	}
	if result = regexp.MustCompile(`^(.+)<(\d+)><([BOT]+)><(.+)>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: result[4]}
	}
	if result = regexp.MustCompile(`^(.+)<(\d+)><([BOT]+)><>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: "Unassigned"}
	}
	if result = regexp.MustCompile(`^(.+)<(\d+)><([BOT]+)>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: "Unassigned"}
	}
	if result = regexp.MustCompile(`^(.+)<(\d+)><([cCoOnNsSoOlLeE]+)>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: "Unassigned"}
	}
	if result = regexp.MustCompile(`^(.+)<(\d+)><([cCoOnNsSoOlLeE]+)><([cCoOnNsSoOlLeE]+)>$`).FindStringSubmatch(playerStr); result != nil {
		id, _ := strconv.Atoi(result[2])
		return &Player{Name: result[1], ID: id, SteamID: result[3], Team: result[4]}
	}

	// Fallback if no pattern matches
	// The original JS returns undefined, we return nil
	return nil
}

// extractIntPosition парсит 3 строковых координаты в срез int
func extractIntPosition(coords []string) ([]int, error) {
	if len(coords) != 3 {
		return nil, fmt.Errorf("expected 3 coordinates, got %d", len(coords))
	}
	pos := make([]int, 3)
	var err error
	pos[0], err = strconv.Atoi(coords[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse x coordinate '%s': %w", coords[0], err)
	}
	pos[1], err = strconv.Atoi(coords[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse y coordinate '%s': %w", coords[1], err)
	}
	pos[2], err = strconv.Atoi(coords[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse z coordinate '%s': %w", coords[2], err)
	}
	return pos, nil
}

// ParseLine парсит одну строку лога и возвращает структуру события или nil, если не удалось распарсить
func ParseLine(line string) interface{} {
	// Предварительная обработка
	// В Go строки immutable, нет нужды в Buffer.IsBuffer
	idx := strings.Index(line, "L ")
	if idx == -1 || len(line) < 30 { // Also check minimum length
		return nil
	}
	processedLine := line[30:]
	processedLine = strings.TrimRight(processedLine, "\r\n\u0000") // Remove \r, \n, \u0000 from end
	processedLine = strings.TrimSpace(processedLine) // Trim leading/trailing whitespace

	// Функции парсинга в Go
	parsers := []func(string) interface{}{
		parseTeamName,
		parseAttacked,
		parseKilledHs,
		parseKilledPenHs,
		parseKilledPen,
		parseKilled,
		parseMolotov,
		parseTriggerTeam,
		parseTriggerWorldOn,
		parseTriggerWorld,
		parseTriggerPlayer,
		parseTriggerPlayerVal,
		parsePurchase,
		parseDisconnected,
		parseEntered,
		parseStartLog,
		parseCloseLog,
		parseConnected,
		parseConnectedAddress,
		parseRcon,
		parseThrew,
		parsePickedUp,
		parseSuicide,
		parseSuicidePosition2,
		parseSuicidePosition,
		parseChangedRole,
		parseSIDValidated,
		parseJoinedTeam,
		parseSwitchTeam,
		parseLoadingMap,
		parseStartedMap,
		parseServerCvar,
		parseScoreCurrent,
		parseScoreScored,
		parseScoreFinal,
		parseSpawned,
		parseSay,
		parseSayTeam,
		parsePositionReport,
		parseAssist,
		parseVoteSuccessKick,
		parseVoteFalseKick,
		parseVoteSuccessNxLvl,
		parseVoteFalseNxLvl,
		parseVoteSuccessChLvl,
		parseVoteFalseChLvl,
		parseVoteSuccessScramble,
		parseVoteFalseScramble,
		parseVoteSuccessSwap,
		parseVoteFalseSwap,
		// hlx2
		parseTriggerKillLocation,
		parseTriggerWeaponstats,
		parseTriggerWeaponstats2,
	}

	var matchedResults []interface{}

	// Последовательная проверка шаблонов
	for _, parserFn := range parsers {
		if result := parserFn(processedLine); result != nil {
			matchedResults = append(matchedResults, result)
		}
	}

	// Обработка результатов
	if len(matchedResults) > 1 {
		// Логирование неправильно разобранных строк
		// В Go лучше использовать логгер или передать функцию обработки ошибок
		log.Printf("Ошибка парсинга: строка '%s' совпала с несколькими шаблонами: %+v\n", processedLine, matchedResults)
		// Можно вернуть nil или первый результат, в зависимости от логики
		return nil // Возвращаем nil для обозначения неопределенного результата
	} else if len(matchedResults) == 1 {
		// Добавляем временную метку
		now := time.Now()
		// Проверяем, поддерживает ли результат интерфейс hasBaseEvent
		if baseEventGetter, ok := matchedResults[0].(hasBaseEvent); ok {
			baseEvent := baseEventGetter.GetBaseEvent()
			baseEvent.Date = EventDate{
				Year:    now.Year(),
				Month:   int(now.Month()), // Convert Month type to int
				Day:     now.Day(),
				Hour:    now.Hour(),
				Minutes: now.Minute(),
				Seconds: now.Second(),
				TS:      now.UnixNano() / int64(time.Millisecond), // Milliseconds timestamp
				FullDate: now,
			}
		} else {
			// Если результат не поддерживает GetBaseEvent, логируем это
			log.Printf("Распарсенное событие не поддерживает GetBaseEvent: %+v", matchedResults[0])
		}

		return matchedResults[0]
	} else {
		// Логирование неразобранных строк
		log.Printf("Строка не распаршена: '%s'\n", processedLine)
		return nil
	}
}

// Вспомогательные функции для каждого типа события

func parseTeamName(line string) interface{} {
	if match := reTeamName.FindStringSubmatch(line); match != nil {
		return &TeamNameEvent{BaseEvent: BaseEvent{Event: "teamName"}, Team: match[1], TeamName: match[2]}
	}
	return nil
}

func parseAttacked(line string) interface{} {
	if match := reAttacked.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		playerB := parsePlayer(match[5])
		if playerA == nil || playerB == nil { // Ensure players were parsed successfully
            return nil // Or log a more specific error
        }
		posA, _ := extractIntPosition(match[2:5])
		posB, _ := extractIntPosition(match[6:9])
		playerA.Position = posA
		playerB.Position = posB // Note: JS code had posB assigned to PlayerA here, fixed to PlayerB
		damage, _ := strconv.Atoi(match[10])
		damageArmor, _ := strconv.Atoi(match[11])
		health, _ := strconv.Atoi(match[12])
		armor, _ := strconv.Atoi(match[13])

		return &AttackedEvent{
			BaseEvent: BaseEvent{Event: "attacked"},
			PlayerA:   playerA,
			PlayerB:   playerB,
			Weapon:  match[9],
			Damage:  damage,
			DamageArmor: damageArmor,
			Health:  health,
			Armor:   armor,
			Hitgroup: match[14],
		}
	}
	return nil
}

func parseKilledHs(line string) interface{} {
	if match := reKilledHs.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		playerB := parsePlayer(match[5])
        if playerA == nil || playerB == nil {
            return nil
        }
		posA, _ := extractIntPosition(match[2:5])
		posB, _ := extractIntPosition(match[6:9])
		playerA.Position = posA
		playerB.Position = posB // Note: JS code had match[8] for posB, should be match[6:9]

		return &KilledEvent{
			BaseEvent: BaseEvent{Event: "killed"},
			PlayerA:   playerA,
			PlayerB:   playerB,
			Weapon:  match[9], // Note: JS code had match[8] for weapon, should be match[9]
			Headshot:  true,
			Penetrated: false,
		}
	}
	return nil
}

func parseKilledPenHs(line string) interface{} {
	if match := reKilledPenHs.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		playerB := parsePlayer(match[5])
        if playerA == nil || playerB == nil {
            return nil
        }
		posA, _ := extractIntPosition(match[2:5])
		posB, _ := extractIntPosition(match[6:9])
		playerA.Position = posA
		playerB.Position = posB // Note: JS code had match[8] for posB, should be match[6:9]

		return &KilledEvent{
			BaseEvent: BaseEvent{Event: "killed"},
			PlayerA:   playerA,
			PlayerB:   playerB,
			Weapon:  match[9], // Note: JS code had match[8] for weapon, should be match[9]
			Headshot:  true,
			Penetrated: true,
		}
	}
	return nil
}

func parseKilledPen(line string) interface{} {
	if match := reKilledPen.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		playerB := parsePlayer(match[5])
        if playerA == nil || playerB == nil {
            return nil
        }
		posA, _ := extractIntPosition(match[2:5])
		posB, _ := extractIntPosition(match[6:9])
		playerA.Position = posA
		playerB.Position = posB // Note: JS code had match[8] for posB, should be match[6:9]

		return &KilledEvent{
			BaseEvent: BaseEvent{Event: "killed"},
			PlayerA:   playerA,
			PlayerB:   playerB,
			Weapon:  match[9], // Note: JS code had match[8] for weapon, should be match[9]
			Headshot:  false,
			Penetrated: true,
		}
	}
	return nil
}

func parseKilled(line string) interface{} {
	if match := reKilled.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		playerB := parsePlayer(match[5])
        if playerA == nil || playerB == nil {
            return nil
        }
		posA, _ := extractIntPosition(match[2:5])
		posB, _ := extractIntPosition(match[6:9])
		playerA.Position = posA
		playerB.Position = posB // Note: JS code had match[8] for posB, should be match[6:9]

		return &KilledEvent{
			BaseEvent: BaseEvent{Event: "killed"},
			PlayerA:   playerA,
			PlayerB:   playerB,
			Weapon:  match[9],
			Headshot:  false,
			Penetrated: false,
		}
	}
	return nil
}

func parseMolotov(line string) interface{} {
	if match := reMolotov.FindStringSubmatch(line); match != nil {
		pos, _ := extractIntPosition(match[1:4])
		vel, _ := extractIntPosition(match[4:7])
		return &MolotovProjectileEvent{BaseEvent: BaseEvent{Event: "molotovProjectile"}, Position: pos, PositionVelocity: vel}
	}
	return nil
}


func parseTriggerTeam(line string) interface{} {
	if match := reTriggerTeam.FindStringSubmatch(line); match != nil {
		ctScore, _ := strconv.Atoi(match[3])
		tScore, _ := strconv.Atoi(match[4])
		return &TriggerEvent{
			BaseEvent: BaseEvent{Event: "trigger"},
			TriggerType: "team",
			Team: match[1],
			Trigger: match[2],
			ResultScoreCT: ctScore,
			ResultScoreT: tScore,
		}
	}
	return nil
}

func parseTriggerWorldOn(line string) interface{} {
	if match := reTriggerWorldOn.FindStringSubmatch(line); match != nil {
		return &TriggerEvent{BaseEvent: BaseEvent{Event: "trigger"}, TriggerType: "world", Trigger: match[1], TriggerOn: match[2]}
	}
	return nil
}

func parseTriggerWorld(line string) interface{} {
	if match := reTriggerWorld.FindStringSubmatch(line); match != nil {
		return &TriggerEvent{BaseEvent: BaseEvent{Event: "trigger"}, TriggerType: "world", Trigger: match[1]}
	}
	return nil
}

func parseTriggerPlayer(line string) interface{} {
	if match := reTriggerPlayer.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &TriggerEvent{BaseEvent: BaseEvent{Event: "trigger"}, TriggerType: "player", Player: player, Trigger: match[2]}
	}
	return nil
}

func parseTriggerPlayerVal(line string) interface{} {
	if match := reTriggerPlayerVal.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &TriggerEvent{BaseEvent: BaseEvent{Event: "trigger"}, TriggerType: "player", Player: player, Trigger: match[2], Value: match[3]}
	}
	return nil
}

func parsePurchase(line string) interface{} {
	if match := rePurchase.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &PurchaseEvent{BaseEvent: BaseEvent{Event: "purchase"}, Player: player, Item: match[2]}
	}
	return nil
}

func parseDisconnected(line string) interface{} {
	if match := reDisconnected.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &DisconnectedEvent{BaseEvent: BaseEvent{Event: "disconnected"}, Player: player, Reason: match[2]}
	}
	return nil
}

func parseEntered(line string) interface{} {
	if match := reEntered.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &EnteredEvent{BaseEvent: BaseEvent{Event: "entered"}, Player: player}
	}
	return nil
}

func parseStartLog(line string) interface{} {
	if reStartLog.MatchString(line) {
		return &LogStartEvent{BaseEvent: BaseEvent{Event: "startLog"}}
	}
	return nil
}

func parseCloseLog(line string) interface{} {
	if reCloseLog.MatchString(line) {
		return &LogCloseEvent{BaseEvent: BaseEvent{Event: "closeLog"}}
	}
	return nil
}

func parseConnected(line string) interface{} {
	if match := reConnected.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &ConnectedEvent{BaseEvent: BaseEvent{Event: "connected"}, Player: player, Address: ""} // Address is empty as per regex
	}
	return nil
}

func parseConnectedAddress(line string) interface{} {
	if match := reConnectedAddress.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &ConnectedEvent{BaseEvent: BaseEvent{Event: "connected"}, Player: player, Address: match[2]}
	}
	return nil
}

func parseRcon(line string) interface{} {
	if match := reRcon.FindStringSubmatch(line); match != nil {
		port, _ := strconv.Atoi(match[2])
		return &RconEvent{BaseEvent: BaseEvent{Event: "rcon"}, Address: match[1], Port: port, Command: match[3]}
	}
	return nil
}

func parseThrew(line string) interface{} {
	if match := reThrew.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		pos, _ := extractIntPosition(match[3:6])
		player.Position = pos
		return &ThrewEvent{BaseEvent: BaseEvent{Event: "threw"}, Player: player, Item: match[2]}
	}
	return nil
}

func parsePickedUp(line string) interface{} {
	if match := rePickedUp.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &PickedUpEvent{BaseEvent: BaseEvent{Event: "pickedUp"}, Player: player, Item: match[2]}
	}
	return nil
}

func parseSuicide(line string) interface{} {
	if match := reSuicide.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &SuicideEvent{BaseEvent: BaseEvent{Event: "suicide"}, Player: player, Item: match[2]}
	}
	return nil
}

func parseSuicidePosition2(line string) interface{} {
	if match := reSuicidePosition2.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		pos, _ := extractIntPosition(match[2:5])
		player.Position = pos
		return &SuicideEvent{BaseEvent: BaseEvent{Event: "suicide"}, Player: player, Item: match[5]}
	}
	return nil
}

func parseSuicidePosition(line string) interface{} {
	if match := reSuicidePosition.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		pos, _ := extractIntPosition(match[3:6])
		player.Position = pos
		return &SuicideEvent{BaseEvent: BaseEvent{Event: "suicide"}, Player: player, Item: match[2]}
	}
	return nil
}


func parseChangedRole(line string) interface{} {
	if match := reChangedRole.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &ChangedRoleEvent{BaseEvent: BaseEvent{Event: "changedRole"}, Player: player, Role: match[2]}
	}
	return nil
}

func parseSIDValidated(line string) interface{} {
	if match := reSIDValidated.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &SIDValidatedEvent{BaseEvent: BaseEvent{Event: "SIDValidated"}, Player: player}
	}
	return nil
}

func parseJoinedTeam(line string) interface{} {
	if match := reJoinedTeam.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &JoinedTeamEvent{BaseEvent: BaseEvent{Event: "joinedTeam"}, Player: player, Team: match[2]}
	}
	return nil
}

func parseSwitchTeam(line string) interface{} {
	if match := reSwitchTeam.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &SwitchTeamEvent{BaseEvent: BaseEvent{Event: "switchTeam"}, Player: player, TeamA: match[2], TeamB: match[3]}
	}
	return nil
}

func parseLoadingMap(line string) interface{} {
	if match := reLoadingMap.FindStringSubmatch(line); match != nil {
		return &LoadingMapEvent{BaseEvent: BaseEvent{Event: "loadingMap"}, Map: match[1]}
	}
	return nil
}

func parseStartedMap(line string) interface{} {
	if match := reStartedMap.FindStringSubmatch(line); match != nil {
		return &StartedMapEvent{BaseEvent: BaseEvent{Event: "startedMap"}, Map: match[1], CRC: match[2]}
	}
	return nil
}

func parseServerCvar(line string) interface{} {
	if match := reServerCvar.FindStringSubmatch(line); match != nil {
		return &ServerCvarEvent{BaseEvent: BaseEvent{Event: "serverCvar"}, Parameter: match[1], Value: match[2]}
	}
	return nil
}

func parseScoreCurrent(line string) interface{} {
	if match := reScoreCurrent.FindStringSubmatch(line); match != nil {
		score, _ := strconv.Atoi(match[2])
		players, _ := strconv.Atoi(match[3])
		return &ScoreEvent{BaseEvent: BaseEvent{Event: "score"}, Team: match[1], ScoreType: "current", Score: score, Players: players}
	}
	return nil
}

func parseScoreScored(line string) interface{} {
	if match := reScoreScored.FindStringSubmatch(line); match != nil {
		score, _ := strconv.Atoi(match[2])
		players, _ := strconv.Atoi(match[3])
		return &ScoreEvent{BaseEvent: BaseEvent{Event: "scored"}, Team: match[1], ScoreType: "scored", Score: score, Players: players}
	}
	return nil
}

func parseScoreFinal(line string) interface{} {
	if match := reScoreFinal.FindStringSubmatch(line); match != nil {
		score, _ := strconv.Atoi(match[2])
		players, _ := strconv.Atoi(match[3])
		return &ScoreEvent{BaseEvent: BaseEvent{Event: "scored"}, Team: match[1], ScoreType: "final", Score: score, Players: players}
	}
	return nil
}


func parseSpawned(line string) interface{} {
	if match := reSpawned.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &SpawnedEvent{BaseEvent: BaseEvent{Event: "spawned"}, Player: player, Spawned: match[2]}
	}
	return nil
}

func parseSay(line string) interface{} {
	if match := reSay.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &SayEvent{BaseEvent: BaseEvent{Event: "say"}, Player: player, Text: match[2]}
	}
	return nil
}

func parseSayTeam(line string) interface{} {
	if match := reSayTeam.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &SayTeamEvent{BaseEvent: BaseEvent{Event: "sayTeam"}, Player: player, Text: match[2]}
	}
	return nil
}

func parsePositionReport(line string) interface{} {
	if match := rePositionReport.FindStringSubmatch(line); match != nil {
		player := parsePlayer(match[1])
        if player == nil { return nil }
		return &PositionReportEvent{BaseEvent: BaseEvent{Event: "positionReport"}, Player: player, Text: match[2]}
	}
	return nil
}

func parseAssist(line string) interface{} {
	if match := reAssist.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		playerB := parsePlayer(match[2])
        if playerA == nil || playerB == nil {
            return nil
        }
		return &AssistEvent{BaseEvent: BaseEvent{Event: "assist"}, PlayerA: playerA, PlayerB: playerB}
	}
	return nil
}

func parseVoteSuccessKick(line string) interface{} {
	if match := reVoteSuccessKick.FindStringSubmatch(line); match != nil {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "kick", VoteResult: "success", Player: match[1]}
	}
	return nil
}

func parseVoteFalseKick(line string) interface{} {
	if match := reVoteFalseKick.FindStringSubmatch(line); match != nil {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "kick", VoteResult: "failed ", Player: match[1]}
	}
	return nil
}

func parseVoteSuccessNxLvl(line string) interface{} {
	if match := reVoteSuccessNxLvl.FindStringSubmatch(line); match != nil {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "nextlevel", VoteResult: "success", Map: match[1]}
	}
	return nil
}

func parseVoteFalseNxLvl(line string) interface{} {
	if match := reVoteFalseNxLvl.FindStringSubmatch(line); match != nil {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "nextlevel", VoteResult: "failed ", Map: match[1]}
	}
	return nil
}

func parseVoteSuccessChLvl(line string) interface{} {
	if match := reVoteSuccessChLvl.FindStringSubmatch(line); match != nil {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "changelevel", VoteResult: "success", Map: match[1]}
	}
	return nil
}

func parseVoteFalseChLvl(line string) interface{} {
	if match := reVoteFalseChLvl.FindStringSubmatch(line); match != nil {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "changelevel", VoteResult: "failed ", Map: match[1]}
	}
	return nil
}

func parseVoteSuccessScramble(line string) interface{} {
	if reVoteSuccessScramble.MatchString(line) {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "scrambleteams", VoteResult: "success"}
	}
	return nil
}

func parseVoteFalseScramble(line string) interface{} {
	if reVoteFalseScramble.MatchString(line) {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "scrambleteams", VoteResult: "failed "}
	}
	return nil
}

func parseVoteSuccessSwap(line string) interface{} {
	if reVoteSuccessSwap.MatchString(line) {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "swapteams", VoteResult: "success"}
	}
	return nil
}

func parseVoteFalseSwap(line string) interface{} {
	if reVoteFalseSwap.MatchString(line) {
		return &VotingEvent{BaseEvent: BaseEvent{Event: "voting"}, VoteType: "swapteams", VoteResult: "failed "}
	}
	return nil
}

// hlx2 parsers
func parseTriggerKillLocation(line string) interface{} {
	if match := reTriggerKillLocation.FindStringSubmatch(line); match != nil {
		playerA := &Player{} // Create empty players to add positions
		playerB := &Player{}
		posA, _ := extractIntPosition(match[2:5])
		posB, _ := extractIntPosition(match[5:8])
		playerA.Position = posA
		playerB.Position = posB

		return &TriggerEvent{
			BaseEvent: BaseEvent{Event: "trigger"},
			TriggerType: "world", // Assuming world trigger based on pattern
			Trigger: match[1],
			PlayerA: playerA,
			PlayerB: playerB,
		}
	}
	return nil
}

func parseTriggerWeaponstats(line string) interface{} {
	if match := reTriggerWeaponstats.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		if playerA == nil { return nil }

		return &TriggerEvent{
			BaseEvent: BaseEvent{Event: "trigger"},
			TriggerType: "Weaponstats", // Preserve typo from JS
			PlayerA: playerA,
			Trigger: match[2], // The trigger name, e.g., "weapon_stats"
			Weapon: match[3],
			Shots: match[4],
			Hits: match[5],
			Kills: match[6],
			Headshots: match[7],
			Tks: match[8],
			Damage: match[9],
			Deaths: match[10],
		}
	}
	return nil
}

func parseTriggerWeaponstats2(line string) interface{} {
	if match := reTriggerWeaponstats2.FindStringSubmatch(line); match != nil {
		playerA := parsePlayer(match[1])
		if playerA == nil { return nil }

		return &TriggerEvent{
			BaseEvent: BaseEvent{Event: "trigger"},
			TriggerType: "Weaponstats2", // Preserve name from JS
			PlayerA: playerA,
			Trigger: match[2], // The trigger name
			Weapon: match[3],
			Head: match[4],
			Chest: match[5],
			Stomach: match[6],
			Leftarm: match[7],
			Rightarm: match[8],
			Leftleg: match[9],
			Rigthleg: match[10], // Preserve typo from JS
		}
	}
	return nil
}


// Интерфейс для получения базового события, чтобы можно было добавить время
type hasBaseEvent interface {
	GetBaseEvent() *BaseEvent
}

// Реализация интерфейса для каждой структуры события
func (e *TeamNameEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *AttackedEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *KilledEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *TriggerEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *PurchaseEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *DisconnectedEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *EnteredEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *LogStartEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *LogCloseEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *ConnectedEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *RconEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *ThrewEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *PickedUpEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *SuicideEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *ChangedRoleEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *SIDValidatedEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *JoinedTeamEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *SwitchTeamEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *LoadingMapEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *StartedMapEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *ServerCvarEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *ScoreEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *SpawnedEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *SayEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *SayTeamEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *PositionReportEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *AssistEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *VotingEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }
func (e *MolotovProjectileEvent) GetBaseEvent() *BaseEvent { return &e.BaseEvent }


// logParseHandler обрабатывает входящие POST запросы с лог-строками
func logParseHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это POST запрос
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		log.Printf("Ошибка чтения тела запроса: %v", err)
		return
	}
	defer r.Body.Close()

	// Ожидаем JSON массив строк в теле запроса
	var logLines []string
	err = json.Unmarshal(body, &logLines)
	if err != nil {
		// Если это не JSON массив, попробуем обработать тело как одну строку или построчно
		// Для простоты сейчас предполагаем JSON массив. Можно добавить обработку текста.
		http.Error(w, "Некорректный формат тела запроса, ожидается JSON массив строк", http.StatusBadRequest)
		log.Printf("Ошибка декодирования JSON: %v", err)
		return
	}

	// Обрабатываем каждую строку лога
	var parsedEvents []interface{}
	for _, line := range logLines {
		event := ParseLine(line)
		if event != nil {
			parsedEvents = append(parsedEvents, event)
		}
	}

	// Отправляем распарсенные события в формате JSON в ответе
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Опционально: для читаемости JSON
	if err := encoder.Encode(parsedEvents); err != nil {
		// Если отправка ответа не удалась, это может быть проблемой сети
		log.Printf("Ошибка отправки JSON ответа: %v", err)
		// Нет смысла вызывать http.Error после того, как уже начали писать в w
		return
	}

	log.Printf("Обработано %d строк лога, найдено %d событий", len(logLines), len(parsedEvents))
}


func main() {
	// Настраиваем маршрутизатор HTTP. Можно использовать gorilla/mux для более сложных маршрутов.
	// Для нашего случая стандартного http.ServeMux достаточно.
	mux := http.NewServeMux()

	// Регистрируем обработчик для POST запросов на пути "/parse"
	mux.HandleFunc("/parse", logParseHandler)

	// Запускаем HTTP сервер на порту 5000
	port := "5000"
	log.Printf("Сервер парсера запущен и слушает на порту %s...", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
