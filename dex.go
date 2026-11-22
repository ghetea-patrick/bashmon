package main

import (
	"fmt"
	"strings"
)

type MonsterType string

type DamageClass uint8

const (
	TypeHardware MonsterType = "Hardware"
	TypeSoftware MonsterType = "Software"
	TypeNetwork  MonsterType = "Network"
	TypeSecurity MonsterType = "Security"
	TypeLogic    MonsterType = "Logic"
	TypeData     MonsterType = "Data"
	TypeMalware  MonsterType = "Malware"
	TypeGlitch   MonsterType = "Glitch"
	TypeMeme     MonsterType = "Meme"
)

const (
	DamagePhysical DamageClass = iota
	DamageSpecial
)

type Move struct {
	Name     string
	Type     MonsterType
	Class    DamageClass
	Power    int
	Accuracy int
	MaxPP    int
	PP       int
	Recoil   bool
}

type Stats struct {
	HP             int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
}

type Monster struct {
	DexNumber int
	Name      string
	Types     [2]MonsterType
	Level     int
	Stats     Stats
	HP        int
	Moves     [4]Move
}

const NumberOfMonsters = 42

var typeChart = map[MonsterType]map[MonsterType]float64{
	TypeHardware: {
		TypeNetwork:  2.0,
		TypeData:     2.0,
		TypeMalware:  0.5,
		TypeSoftware: 0.5,
	},
	TypeSoftware: {
		TypeHardware: 2.0,
		TypeLogic:    2.0,
		TypeMalware:  0.5,
		TypeSecurity: 0.5,
	},
	TypeNetwork: {
		TypeSoftware: 2.0,
		TypeMeme:     2.0,
		TypeHardware: 0.5,
		TypeGlitch:   0.5,
	},
	TypeSecurity: {
		TypeMalware: 2.0,
		TypeGlitch:  2.0,
		TypeLogic:   0.5,
		TypeMeme:    0.5,
	},
	TypeLogic: {
		TypeMeme:     2.0,
		TypeMalware:  2.0,
		TypeSoftware: 0.5,
		TypeData:     0.5,
	},
	TypeData: {
		TypeLogic:    2.0,
		TypeSecurity: 2.0,
		TypeHardware: 0.5,
		TypeMalware:  0.5,
	},
	TypeMalware: {
		TypeSoftware: 2.0,
		TypeData:     2.0,
		TypeSecurity: 0.5,
		TypeLogic:    0.5,
	},
	TypeGlitch: {
		TypeHardware: 2.0,
		TypeNetwork:  2.0,
		TypeSecurity: 0.5,
		TypeLogic:    0.5,
	},
	TypeMeme: {
		TypeSecurity: 2.0,
		TypeLogic:    2.0,
		TypeNetwork:  0.5,
		TypeMalware:  0.5,
	},
}

var KernelPanic = Move{
	Name:     "Kernel Panic",
	Type:     TypeGlitch,
	Class:    DamagePhysical,
	Power:    35,
	Accuracy: 100,
	Recoil:   true,
}

var MonsterDex = [...]Monster{
	makeSpecies(1, "Bytebeast", [2]MonsterType{TypeData, TypeHardware}, Stats{HP: 44, Attack: 29, Defense: 27, SpecialAttack: 27, SpecialDefense: 25, Speed: 23},
		makeMove("Byte Bash", TypeData, DamagePhysical, 40, 100, 25),
		makeMove("Cache Crash", TypeHardware, DamagePhysical, 55, 95, 15),
		makeMove("Checksum", TypeData, DamageSpecial, 45, 100, 20),
		makeMove("Static Spark", TypeGlitch, DamageSpecial, 60, 90, 10)),
	makeSpecies(2, "Cachecow", [2]MonsterType{TypeData, TypeHardware}, Stats{HP: 54, Attack: 27, Defense: 33, SpecialAttack: 25, SpecialDefense: 31, Speed: 16},
		makeMove("Cache Slap", TypeData, DamagePhysical, 45, 100, 25),
		makeMove("Disk Stampede", TypeHardware, DamagePhysical, 55, 95, 15),
		makeMove("Buffer Moo", TypeSoftware, DamageSpecial, 50, 95, 15),
		makeMove("Cold Storage", TypeSecurity, DamageSpecial, 40, 100, 20)),
	makeSpecies(3, "Stackrat", [2]MonsterType{TypeLogic, TypeData}, Stats{HP: 39, Attack: 31, Defense: 22, SpecialAttack: 28, SpecialDefense: 24, Speed: 32},
		makeMove("Stack Jab", TypeLogic, DamagePhysical, 40, 100, 25),
		makeMove("Overflow", TypeData, DamageSpecial, 60, 90, 10),
		makeMove("Last In", TypeLogic, DamageSpecial, 45, 100, 20),
		makeMove("Rat Race", TypeNetwork, DamagePhysical, 50, 95, 15)),
	makeSpecies(4, "Heapling", [2]MonsterType{TypeHardware, TypeData}, Stats{HP: 48, Attack: 25, Defense: 33, SpecialAttack: 26, SpecialDefense: 29, Speed: 19},
		makeMove("Heap Hop", TypeHardware, DamagePhysical, 45, 100, 20),
		makeMove("Memory Leak", TypeData, DamageSpecial, 55, 95, 15),
		makeMove("Garbage Toss", TypeSoftware, DamagePhysical, 50, 95, 15),
		makeMove("Free Block", TypeSecurity, DamageSpecial, 40, 100, 20)),
	makeSpecies(5, "Nullpuff", [2]MonsterType{TypeGlitch, TypeLogic}, Stats{HP: 42, Attack: 21, Defense: 27, SpecialAttack: 34, SpecialDefense: 29, Speed: 27},
		makeMove("Null Nudge", TypeGlitch, DamageSpecial, 40, 100, 25),
		makeMove("Void Pop", TypeGlitch, DamageSpecial, 60, 90, 10),
		makeMove("False Return", TypeLogic, DamageSpecial, 50, 95, 15),
		makeMove("Crash Cushion", TypeHardware, DamagePhysical, 45, 100, 20)),
	makeSpecies(6, "Segfault", [2]MonsterType{TypeSoftware, TypeGlitch}, Stats{HP: 40, Attack: 24, Defense: 22, SpecialAttack: 36, SpecialDefense: 25, Speed: 31},
		makeMove("Bad Access", TypeGlitch, DamageSpecial, 55, 95, 15),
		makeMove("Pointer Poke", TypeData, DamagePhysical, 45, 100, 20),
		makeMove("Core Dump", TypeSoftware, DamageSpecial, 65, 85, 10),
		makeMove("Try Again", TypeLogic, DamageSpecial, 40, 100, 25)),
	makeSpecies(7, "Bugaboo", [2]MonsterType{TypeGlitch, TypeMalware}, Stats{HP: 43, Attack: 31, Defense: 23, SpecialAttack: 28, SpecialDefense: 24, Speed: 28},
		makeMove("Glitch Bite", TypeGlitch, DamagePhysical, 40, 100, 25),
		makeMove("Bug Spray", TypeSecurity, DamageSpecial, 50, 95, 15),
		makeMove("Crash Crawl", TypeMalware, DamagePhysical, 55, 95, 15),
		makeMove("Patch Panic", TypeSoftware, DamageSpecial, 60, 90, 10)),
	makeSpecies(8, "Glitchimp", [2]MonsterType{TypeGlitch, TypeMeme}, Stats{HP: 41, Attack: 29, Defense: 22, SpecialAttack: 33, SpecialDefense: 23, Speed: 34},
		makeMove("Pixel Prank", TypeMeme, DamageSpecial, 45, 100, 20),
		makeMove("Reality Lag", TypeGlitch, DamageSpecial, 60, 90, 10),
		makeMove("Banana Buffer", TypeData, DamagePhysical, 50, 95, 15),
		makeMove("Oopsie Loop", TypeLogic, DamageSpecial, 40, 100, 25)),
	makeSpecies(9, "Loopaloop", [2]MonsterType{TypeLogic, TypeSoftware}, Stats{HP: 45, Attack: 27, Defense: 26, SpecialAttack: 32, SpecialDefense: 28, Speed: 27},
		makeMove("While Whip", TypeLogic, DamagePhysical, 45, 100, 20),
		makeMove("Endless Cycle", TypeSoftware, DamageSpecial, 55, 95, 15),
		makeMove("Do Again", TypeLogic, DamageSpecial, 40, 100, 25),
		makeMove("Loop Crash", TypeGlitch, DamageSpecial, 60, 90, 10)),
	makeSpecies(10, "Forknado", [2]MonsterType{TypeLogic, TypeNetwork}, Stats{HP: 42, Attack: 33, Defense: 24, SpecialAttack: 28, SpecialDefense: 25, Speed: 33},
		makeMove("Fork Strike", TypeLogic, DamagePhysical, 50, 100, 20),
		makeMove("Thread Storm", TypeNetwork, DamageSpecial, 60, 90, 10),
		makeMove("Merge Clash", TypeSoftware, DamagePhysical, 55, 95, 15),
		makeMove("Split Call", TypeLogic, DamageSpecial, 45, 100, 20)),
	makeSpecies(11, "Recursnake", [2]MonsterType{TypeLogic, TypeGlitch}, Stats{HP: 44, Attack: 31, Defense: 25, SpecialAttack: 30, SpecialDefense: 26, Speed: 32},
		makeMove("Tail Call", TypeLogic, DamagePhysical, 45, 100, 20),
		makeMove("Stack Squeeze", TypeData, DamageSpecial, 60, 90, 10),
		makeMove("Base Case", TypeLogic, DamageSpecial, 40, 100, 25),
		makeMove("Infinite Hiss", TypeGlitch, DamageSpecial, 55, 95, 15)),
	makeSpecies(12, "Booloon", [2]MonsterType{TypeLogic, TypeMeme}, Stats{HP: 46, Attack: 23, Defense: 25, SpecialAttack: 34, SpecialDefense: 31, Speed: 24},
		makeMove("True Pop", TypeLogic, DamageSpecial, 45, 100, 20),
		makeMove("False Alarm", TypeMeme, DamageSpecial, 55, 95, 15),
		makeMove("Maybe Missile", TypeSoftware, DamageSpecial, 60, 90, 10),
		makeMove("If Bonk", TypeLogic, DamagePhysical, 45, 100, 20)),
	makeSpecies(13, "Bitrot", [2]MonsterType{TypeMalware, TypeData}, Stats{HP: 46, Attack: 30, Defense: 24, SpecialAttack: 34, SpecialDefense: 27, Speed: 22},
		makeMove("Rotten Byte", TypeMalware, DamagePhysical, 50, 100, 20),
		makeMove("Data Decay", TypeData, DamageSpecial, 55, 95, 15),
		makeMove("Corrupt File", TypeMalware, DamageSpecial, 65, 85, 10),
		makeMove("Checksum Chew", TypeHardware, DamagePhysical, 40, 100, 25)),
	makeSpecies(14, "Packetyak", [2]MonsterType{TypeNetwork, TypeHardware}, Stats{HP: 51, Attack: 31, Defense: 32, SpecialAttack: 24, SpecialDefense: 27, Speed: 18},
		makeMove("Packet Ram", TypeNetwork, DamagePhysical, 45, 100, 20),
		makeMove("Payload Punt", TypeData, DamagePhysical, 55, 95, 15),
		makeMove("Route Rush", TypeNetwork, DamageSpecial, 60, 90, 10),
		makeMove("Ether Bonk", TypeHardware, DamagePhysical, 50, 100, 20)),
	makeSpecies(15, "Pingwing", [2]MonsterType{TypeNetwork, TypeGlitch}, Stats{HP: 40, Attack: 25, Defense: 23, SpecialAttack: 32, SpecialDefense: 25, Speed: 36},
		makeMove("Ping Peck", TypeNetwork, DamagePhysical, 40, 100, 25),
		makeMove("Echo Dive", TypeGlitch, DamageSpecial, 55, 95, 15),
		makeMove("Timeout Hit", TypeNetwork, DamageSpecial, 60, 90, 10),
		makeMove("Packet Flap", TypeHardware, DamagePhysical, 45, 100, 20)),
	makeSpecies(16, "Lagomorph", [2]MonsterType{TypeNetwork, TypeLogic}, Stats{HP: 43, Attack: 27, Defense: 24, SpecialAttack: 29, SpecialDefense: 26, Speed: 38},
		makeMove("Delay Hop", TypeNetwork, DamagePhysical, 45, 100, 20),
		makeMove("Queue Jump", TypeLogic, DamagePhysical, 50, 100, 20),
		makeMove("Buffer Burst", TypeData, DamageSpecial, 60, 90, 10),
		makeMove("Fast Forward", TypeSoftware, DamageSpecial, 55, 95, 15)),
	makeSpecies(17, "Routertoad", [2]MonsterType{TypeNetwork, TypeHardware}, Stats{HP: 49, Attack: 31, Defense: 31, SpecialAttack: 27, SpecialDefense: 28, Speed: 20},
		makeMove("Port Slam", TypeHardware, DamagePhysical, 50, 100, 20),
		makeMove("Route Reboot", TypeNetwork, DamageSpecial, 55, 95, 15),
		makeMove("Hop Count", TypeNetwork, DamageSpecial, 45, 100, 20),
		makeMove("Firewall", TypeSecurity, DamagePhysical, 60, 90, 10)),
	makeSpecies(18, "Sockpuppet", [2]MonsterType{TypeMalware, TypeNetwork}, Stats{HP: 44, Attack: 29, Defense: 25, SpecialAttack: 33, SpecialDefense: 26, Speed: 29},
		makeMove("Proxy Punch", TypeNetwork, DamagePhysical, 45, 100, 20),
		makeMove("Maskerade", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Fake Reply", TypeMalware, DamageSpecial, 50, 95, 15),
		makeMove("Packet Snare", TypeNetwork, DamagePhysical, 55, 95, 15)),
	makeSpecies(19, "Firewally", [2]MonsterType{TypeSecurity, TypeHardware}, Stats{HP: 52, Attack: 30, Defense: 35, SpecialAttack: 25, SpecialDefense: 33, Speed: 17},
		makeMove("Rulebook Bash", TypeSecurity, DamagePhysical, 45, 100, 20),
		makeMove("Port Shield", TypeHardware, DamageSpecial, 55, 95, 15),
		makeMove("Block Party", TypeSecurity, DamageSpecial, 60, 90, 10),
		makeMove("Steel Policy", TypeHardware, DamagePhysical, 50, 100, 20)),
	makeSpecies(20, "Portcorgi", [2]MonsterType{TypeNetwork, TypeSecurity}, Stats{HP: 45, Attack: 29, Defense: 28, SpecialAttack: 29, SpecialDefense: 30, Speed: 29},
		makeMove("Port Nibble", TypeNetwork, DamagePhysical, 40, 100, 25),
		makeMove("Guard Bark", TypeSecurity, DamageSpecial, 55, 95, 15),
		makeMove("Socket Sprint", TypeNetwork, DamagePhysical, 50, 100, 20),
		makeMove("Access Denied", TypeSecurity, DamageSpecial, 60, 90, 10)),
	makeSpecies(21, "Handshark", [2]MonsterType{TypeNetwork, TypeMalware}, Stats{HP: 48, Attack: 34, Defense: 26, SpecialAttack: 29, SpecialDefense: 26, Speed: 28},
		makeMove("Syn Bite", TypeNetwork, DamagePhysical, 45, 100, 20),
		makeMove("Handshake", TypeNetwork, DamagePhysical, 55, 95, 15),
		makeMove("Trojan Fin", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Three Way", TypeSecurity, DamageSpecial, 50, 95, 15)),
	makeSpecies(22, "Namegnome", [2]MonsterType{TypeData, TypeLogic}, Stats{HP: 43, Attack: 23, Defense: 27, SpecialAttack: 34, SpecialDefense: 33, Speed: 21},
		makeMove("Name Tag", TypeData, DamageSpecial, 45, 100, 20),
		makeMove("Symbol Swap", TypeLogic, DamageSpecial, 55, 95, 15),
		makeMove("Scope Snare", TypeSoftware, DamageSpecial, 60, 90, 10),
		makeMove("Root Word", TypeData, DamagePhysical, 50, 100, 20)),
	makeSpecies(23, "Cookiemon", [2]MonsterType{TypeSoftware, TypeMeme}, Stats{HP: 47, Attack: 26, Defense: 25, SpecialAttack: 33, SpecialDefense: 28, Speed: 27},
		makeMove("Cookie Crush", TypeSoftware, DamageSpecial, 45, 100, 20),
		makeMove("Session Snack", TypeData, DamageSpecial, 55, 95, 15),
		makeMove("Pop-up Prank", TypeMeme, DamageSpecial, 60, 90, 10),
		makeMove("Cache Treat", TypeHardware, DamagePhysical, 40, 100, 25)),
	makeSpecies(24, "Phishstick", [2]MonsterType{TypeMalware, TypeNetwork}, Stats{HP: 42, Attack: 27, Defense: 23, SpecialAttack: 36, SpecialDefense: 25, Speed: 31},
		makeMove("Urgent Click", TypeMalware, DamageSpecial, 55, 95, 15),
		makeMove("Hook Line", TypeNetwork, DamagePhysical, 45, 100, 20),
		makeMove("Fake Login", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Spam Jab", TypeMeme, DamagePhysical, 50, 95, 15)),
	makeSpecies(25, "Spamurai", [2]MonsterType{TypeMalware, TypeMeme}, Stats{HP: 45, Attack: 34, Defense: 25, SpecialAttack: 29, SpecialDefense: 25, Speed: 30},
		makeMove("Spam Slash", TypeMalware, DamagePhysical, 50, 100, 20),
		makeMove("Reply All", TypeNetwork, DamageSpecial, 45, 100, 20),
		makeMove("Caps Lock", TypeMeme, DamagePhysical, 55, 95, 15),
		makeMove("Inbox Storm", TypeMalware, DamageSpecial, 60, 90, 10)),
	makeSpecies(26, "Trollbyte", [2]MonsterType{TypeMeme, TypeMalware}, Stats{HP: 46, Attack: 32, Defense: 25, SpecialAttack: 30, SpecialDefense: 27, Speed: 29},
		makeMove("Ratio Rake", TypeMeme, DamagePhysical, 50, 100, 20),
		makeMove("Flamebait", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Skill Issue", TypeMeme, DamageSpecial, 55, 95, 15),
		makeMove("Reply Guy", TypeNetwork, DamagePhysical, 45, 100, 20)),
	makeSpecies(27, "Botfly", [2]MonsterType{TypeSoftware, TypeMalware}, Stats{HP: 44, Attack: 29, Defense: 24, SpecialAttack: 33, SpecialDefense: 25, Speed: 32},
		makeMove("Auto Click", TypeSoftware, DamagePhysical, 45, 100, 20),
		makeMove("Script Swarm", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Macro Maul", TypeSoftware, DamagePhysical, 55, 95, 15),
		makeMove("Human Check", TypeSecurity, DamageSpecial, 50, 95, 15)),
	makeSpecies(28, "Hashbrown", [2]MonsterType{TypeData, TypeSecurity}, Stats{HP: 48, Attack: 30, Defense: 31, SpecialAttack: 29, SpecialDefense: 32, Speed: 20},
		makeMove("Salt Sprinkle", TypeData, DamageSpecial, 45, 100, 20),
		makeMove("Hash Smash", TypeData, DamagePhysical, 55, 95, 15),
		makeMove("Pepper Shield", TypeSecurity, DamageSpecial, 60, 90, 10),
		makeMove("Brute Force", TypeLogic, DamagePhysical, 50, 100, 20)),
	makeSpecies(29, "Rootkitten", [2]MonsterType{TypeMalware, TypeGlitch}, Stats{HP: 42, Attack: 33, Defense: 23, SpecialAttack: 34, SpecialDefense: 24, Speed: 33},
		makeMove("Kernel Claw", TypeMalware, DamagePhysical, 50, 100, 20),
		makeMove("Hidden Proc", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Privilege", TypeSecurity, DamagePhysical, 55, 95, 15),
		makeMove("Ghost Driver", TypeGlitch, DamageSpecial, 50, 95, 15)),
	makeSpecies(30, "Ransomclam", [2]MonsterType{TypeMalware, TypeSecurity}, Stats{HP: 53, Attack: 31, Defense: 34, SpecialAttack: 31, SpecialDefense: 34, Speed: 15},
		makeMove("Encrypt Snap", TypeMalware, DamagePhysical, 55, 95, 15),
		makeMove("Key Demand", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Backup Pinch", TypeSecurity, DamageSpecial, 45, 100, 20),
		makeMove("Restore Shell", TypeSoftware, DamageSpecial, 50, 95, 15)),
	makeSpecies(31, "Clockroach", [2]MonsterType{TypeHardware, TypeLogic}, Stats{HP: 46, Attack: 31, Defense: 31, SpecialAttack: 25, SpecialDefense: 29, Speed: 26},
		makeMove("Tick Tackle", TypeHardware, DamagePhysical, 45, 100, 20),
		makeMove("Clock Drift", TypeLogic, DamageSpecial, 60, 90, 10),
		makeMove("Time Slice", TypeSoftware, DamageSpecial, 55, 95, 15),
		makeMove("Gear Grind", TypeHardware, DamagePhysical, 50, 100, 20)),
	makeSpecies(32, "Pixelhog", [2]MonsterType{TypeData, TypeHardware}, Stats{HP: 49, Attack: 32, Defense: 29, SpecialAttack: 29, SpecialDefense: 27, Speed: 22},
		makeMove("Pixel Pound", TypeData, DamagePhysical, 50, 100, 20),
		makeMove("Render Rush", TypeSoftware, DamageSpecial, 60, 90, 10),
		makeMove("Color Crash", TypeGlitch, DamageSpecial, 55, 95, 15),
		makeMove("Screen Slam", TypeHardware, DamagePhysical, 45, 100, 20)),
	makeSpecies(33, "Ramalam", [2]MonsterType{TypeHardware, TypeSoftware}, Stats{HP: 47, Attack: 31, Defense: 29, SpecialAttack: 32, SpecialDefense: 29, Speed: 25},
		makeMove("RAM Jam", TypeHardware, DamagePhysical, 55, 95, 15),
		makeMove("Cache Cannon", TypeData, DamageSpecial, 60, 90, 10),
		makeMove("Swap File", TypeSoftware, DamageSpecial, 50, 95, 15),
		makeMove("Memory Wall", TypeSecurity, DamagePhysical, 45, 100, 20)),
	makeSpecies(34, "Floppup", [2]MonsterType{TypeHardware, TypeData}, Stats{HP: 43, Attack: 27, Defense: 30, SpecialAttack: 27, SpecialDefense: 30, Speed: 22},
		makeMove("Disk Flick", TypeHardware, DamagePhysical, 45, 100, 20),
		makeMove("Read Error", TypeGlitch, DamageSpecial, 55, 95, 15),
		makeMove("Write Protect", TypeSecurity, DamageSpecial, 60, 90, 10),
		makeMove("Retro Spin", TypeData, DamagePhysical, 50, 100, 20)),
	makeSpecies(35, "Disko", [2]MonsterType{TypeHardware, TypeNetwork}, Stats{HP: 48, Attack: 33, Defense: 31, SpecialAttack: 25, SpecialDefense: 26, Speed: 23},
		makeMove("Spin Cycle", TypeHardware, DamagePhysical, 50, 100, 20),
		makeMove("Seek & Destroy", TypeNetwork, DamagePhysical, 55, 95, 15),
		makeMove("Bad Sector", TypeGlitch, DamageSpecial, 60, 90, 10),
		makeMove("Drive By", TypeHardware, DamagePhysical, 45, 100, 20)),
	makeSpecies(36, "Bootleg", [2]MonsterType{TypeSoftware, TypeGlitch}, Stats{HP: 43, Attack: 27, Defense: 24, SpecialAttack: 35, SpecialDefense: 28, Speed: 29},
		makeMove("Fake Boot", TypeSoftware, DamageSpecial, 50, 100, 20),
		makeMove("Blue Screen", TypeGlitch, DamageSpecial, 65, 85, 10),
		makeMove("Patchwork", TypeSoftware, DamagePhysical, 45, 100, 20),
		makeMove("Boot Loop", TypeLogic, DamageSpecial, 55, 95, 15)),
	makeSpecies(37, "Wirewraith", [2]MonsterType{TypeNetwork, TypeGlitch}, Stats{HP: 42, Attack: 28, Defense: 24, SpecialAttack: 35, SpecialDefense: 27, Speed: 33},
		makeMove("Ghost Packet", TypeNetwork, DamageSpecial, 55, 95, 15),
		makeMove("Static Haunt", TypeGlitch, DamageSpecial, 60, 90, 10),
		makeMove("Cable Lash", TypeHardware, DamagePhysical, 45, 100, 20),
		makeMove("Signal Fade", TypeNetwork, DamageSpecial, 50, 95, 15)),
	makeSpecies(38, "Qwertybug", [2]MonsterType{TypeHardware, TypeMeme}, Stats{HP: 45, Attack: 31, Defense: 27, SpecialAttack: 28, SpecialDefense: 26, Speed: 31},
		makeMove("Key Mash", TypeMeme, DamagePhysical, 45, 100, 20),
		makeMove("Shift Smash", TypeHardware, DamagePhysical, 55, 95, 15),
		makeMove("Typo Trap", TypeLogic, DamageSpecial, 60, 90, 10),
		makeMove("Caps Lock", TypeMeme, DamageSpecial, 50, 100, 20)),
	makeSpecies(39, "Todozer", [2]MonsterType{TypeLogic, TypeSoftware}, Stats{HP: 50, Attack: 32, Defense: 29, SpecialAttack: 28, SpecialDefense: 30, Speed: 20},
		makeMove("Task Stamp", TypeLogic, DamagePhysical, 50, 100, 20),
		makeMove("Deadline Dash", TypeSoftware, DamagePhysical, 55, 95, 15),
		makeMove("Checklist", TypeSecurity, DamageSpecial, 45, 100, 20),
		makeMove("Scope Creep", TypeMeme, DamageSpecial, 60, 90, 10)),
	makeSpecies(40, "Gitgoblin", [2]MonsterType{TypeSoftware, TypeData}, Stats{HP: 44, Attack: 30, Defense: 25, SpecialAttack: 33, SpecialDefense: 27, Speed: 30},
		makeMove("Commit Crush", TypeSoftware, DamagePhysical, 50, 100, 20),
		makeMove("Conflict", TypeLogic, DamageSpecial, 60, 90, 10),
		makeMove("Branch Bash", TypeData, DamagePhysical, 45, 100, 20),
		makeMove("Rebase Rage", TypeSoftware, DamageSpecial, 55, 95, 15)),
	makeSpecies(41, "Memeory", [2]MonsterType{TypeMeme, TypeData}, Stats{HP: 46, Attack: 25, Defense: 27, SpecialAttack: 34, SpecialDefense: 31, Speed: 26},
		makeMove("Deep Fried", TypeMeme, DamageSpecial, 50, 100, 20),
		makeMove("Recall Roast", TypeData, DamageSpecial, 55, 95, 15),
		makeMove("Brainrot", TypeMalware, DamageSpecial, 60, 90, 10),
		makeMove("Archive Flex", TypeSecurity, DamageSpecial, 45, 100, 20)),
	makeSpecies(42, "Yeetabyte", [2]MonsterType{TypeData, TypeGlitch}, Stats{HP: 45, Attack: 31, Defense: 25, SpecialAttack: 36, SpecialDefense: 26, Speed: 31},
		makeMove("Mass Upload", TypeData, DamageSpecial, 55, 95, 15),
		makeMove("Yeet Stream", TypeGlitch, DamageSpecial, 60, 90, 10),
		makeMove("Byte Blast", TypeHardware, DamagePhysical, 50, 100, 20),
		makeMove("Quantum Toss", TypeLogic, DamageSpecial, 65, 85, 10)),
}

func makeSpecies(number int, name string, monsterTypes [2]MonsterType, stats Stats, moves ...Move) Monster {
	var moveSet [4]Move
	copy(moveSet[:], moves)
	return Monster{
		DexNumber: number,
		Name:      name,
		Types:     monsterTypes,
		Level:     12,
		Stats:     stats,
		HP:        stats.HP,
		Moves:     moveSet,
	}
}

func makeMove(name string, moveType MonsterType, class DamageClass, power, accuracy, maxPP int) Move {
	return Move{
		Name:     name,
		Type:     moveType,
		Class:    class,
		Power:    power,
		Accuracy: accuracy,
		MaxPP:    maxPP,
		PP:       maxPP,
	}
}

func freshMonster(template Monster) Monster {
	monster := template
	monster.HP = monster.Stats.HP
	for index := range monster.Moves {
		monster.Moves[index].PP = monster.Moves[index].MaxPP
	}
	return monster
}

func validateDex() error {
	if len(MonsterDex) != NumberOfMonsters {
		return fmt.Errorf("expected %d monsters, found %d", NumberOfMonsters, len(MonsterDex))
	}

	knownTypes := map[MonsterType]bool{
		TypeHardware: true,
		TypeSoftware: true,
		TypeNetwork:  true,
		TypeSecurity: true,
		TypeLogic:    true,
		TypeData:     true,
		TypeMalware:  true,
		TypeGlitch:   true,
		TypeMeme:     true,
	}
	seenNames := make(map[string]bool, len(MonsterDex))

	for index, monster := range MonsterDex {
		if monster.DexNumber != index+1 {
			return fmt.Errorf("dex entry %d has number %d", index+1, monster.DexNumber)
		}
		if strings.TrimSpace(monster.Name) == "" {
			return fmt.Errorf("dex entry %d has an empty name", index+1)
		}
		nameKey := strings.ToLower(monster.Name)
		if seenNames[nameKey] {
			return fmt.Errorf("duplicate monster name %q", monster.Name)
		}
		seenNames[nameKey] = true

		if monster.Level < 1 || monster.Stats.HP < 1 || monster.HP < 1 || monster.HP > monster.Stats.HP {
			return fmt.Errorf("%s has invalid level or health values", monster.Name)
		}
		if monster.Stats.Attack < 1 || monster.Stats.Defense < 1 || monster.Stats.SpecialAttack < 1 || monster.Stats.SpecialDefense < 1 || monster.Stats.Speed < 1 {
			return fmt.Errorf("%s has an invalid combat stat", monster.Name)
		}
		if !knownTypes[monster.Types[0]] {
			return fmt.Errorf("%s has unknown primary type %q", monster.Name, monster.Types[0])
		}
		if monster.Types[1] != "" && (!knownTypes[monster.Types[1]] || monster.Types[1] == monster.Types[0]) {
			return fmt.Errorf("%s has invalid secondary type %q", monster.Name, monster.Types[1])
		}

		for moveIndex, move := range monster.Moves {
			if move.Name == "" {
				return fmt.Errorf("%s has an empty move in slot %d", monster.Name, moveIndex+1)
			}
			if !knownTypes[move.Type] || move.Power < 1 || move.Accuracy < 1 || move.Accuracy > 100 || move.MaxPP < 1 || move.PP != move.MaxPP {
				return fmt.Errorf("%s has invalid move data for %q", monster.Name, move.Name)
			}
			if move.Class != DamagePhysical && move.Class != DamageSpecial {
				return fmt.Errorf("%s has an invalid damage class for %q", monster.Name, move.Name)
			}
		}
	}

	for attackType, matchups := range typeChart {
		if !knownTypes[attackType] {
			return fmt.Errorf("type chart contains unknown attacking type %q", attackType)
		}
		for defenseType, modifier := range matchups {
			if !knownTypes[defenseType] {
				return fmt.Errorf("type chart contains unknown defending type %q", defenseType)
			}
			if modifier != 0.5 && modifier != 2.0 {
				return fmt.Errorf("type chart has unsupported modifier %v for %s against %s", modifier, attackType, defenseType)
			}
		}
	}

	if KernelPanic.Name == "" || KernelPanic.Power < 1 || KernelPanic.Accuracy != 100 || !KernelPanic.Recoil {
		return fmt.Errorf("fallback move is not configured correctly")
	}

	return nil
}
