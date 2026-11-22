package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

const (
	screenWidth        = 34
	screenHeight       = 8
	dexPageSize        = 4
	promptOptionWidth  = 24
	consoleButtonStart = 24
)

const (
	ansiReset      = "\033[0m"
	ansiDim        = "\033[2m"
	ansiHidden     = "\033[8m"
	ansiDPadBlack  = "\033[30m"
	ansiButtonPink = "\033[38;2;255;105;180m"
)

type waitDuration time.Duration

type game struct {
	input  *bufio.Reader
	output io.Writer
}

var errQuit = errors.New("quit requested")

const (
	waitNone   waitDuration = 0
	waitShort  waitDuration = waitDuration(750 * time.Millisecond)
	waitMedium waitDuration = waitDuration(1500 * time.Millisecond)
	waitLong   waitDuration = waitDuration(3 * time.Second)
)

func newGame(input io.Reader, output io.Writer) *game {
	return &game{
		input:  bufio.NewReader(input),
		output: output,
	}
}

func main() {
	if err := validateDex(); err != nil {
		fmt.Fprintf(os.Stderr, "Bashmon data error: %v\n", err)
		os.Exit(1)
	}

	g := newGame(os.Stdin, os.Stdout)
	defer g.resetTerminal()

	signals := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	defer close(done)

	go func() {
		select {
		case <-signals:
			g.resetTerminal()
			os.Exit(130)
		case <-done:
		}
	}()

	if err := g.run(); err != nil && !errors.Is(err, errQuit) && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "Bashmon: %v\n", err)
		os.Exit(1)
	}
}

func (g *game) run() error {
	for {
		g.renderGameBoy(titleScreen(), waitLong)

		playerTemplate, err := g.chooseMonster()
		if err != nil {
			return err
		}

		player := freshMonster(playerTemplate)
		confirmed, err := g.confirmMonster(player)
		if err != nil {
			return err
		}
		if !confirmed {
			continue
		}

		rival := freshMonster(g.chooseRival(player.DexNumber))
		g.showRivalArrival(rival)

		playerWon, err := g.battle(&player, &rival)
		if err != nil {
			return err
		}

		if playerWon {
			g.showVictory(player)
		} else {
			g.showDefeat(player)
		}
	}
}

func (g *game) chooseMonster() (Monster, error) {
	page := 0
	pageCount := (len(MonsterDex) + dexPageSize - 1) / dexPageSize

	for {
		start := page * dexPageSize
		end := min(start+dexPageSize, len(MonsterDex))
		screen := blankScreen()
		for index := start; index < end; index++ {
			slot := index - start + 1
			entry := fmt.Sprintf("[%d] #%02d  %s", slot, MonsterDex[index].DexNumber, MonsterDex[index].Name)
			screen[slot+1] = fitLine(entry)
		}

		g.renderGameBoy(screen, waitNone)
		g.writePrompt("1, 2, 3, 4, Y, N, Q")

		key, err := g.readKey()
		if err != nil {
			return Monster{}, err
		}

		switch key {
		case 'q':
			return Monster{}, errQuit
		case 'y':
			if page+1 < pageCount {
				page++
			}
		case 'n':
			if page > 0 {
				page--
			}
		case '1', '2', '3', '4':
			selected := start + int(key-'1')
			if selected < end {
				return MonsterDex[selected], nil
			}
		}
	}
}

func (g *game) confirmMonster(monster Monster) (bool, error) {
	screen := blankScreen()
	screen[2] = centerLine("So! You want the")
	screen[3] = centerLine(formatMonsterTypes(monster))
	screen[4] = centerLine(monster.Name + "?")
	g.renderGameBoy(screen, waitNone)

	for {
		g.writePrompt("Y, N, Q")
		key, err := g.readKey()
		if err != nil {
			return false, err
		}
		switch key {
		case 'y':
			return true, nil
		case 'n':
			return false, nil
		case 'q':
			return false, errQuit
		}
	}
}

func (g *game) chooseRival(playerDexNumber int) Monster {
	available := make([]Monster, 0, len(MonsterDex)-1)
	for _, candidate := range MonsterDex {
		if candidate.DexNumber != playerDexNumber {
			available = append(available, candidate)
		}
	}
	return available[rand.IntN(len(available))]
}

func (g *game) showRivalArrival(rival Monster) {
	screen := blankScreen()
	screen[3] = centerLine(rival.Name + " wants to fight!")
	screen[4] = centerLine("A wild challenger appeared!")
	g.renderGameBoy(screen, waitLong)
}

func (g *game) battle(player, rival *Monster) (bool, error) {
	playerActsFirst := player.Stats.Speed > rival.Stats.Speed
	if player.Stats.Speed == rival.Stats.Speed {
		playerActsFirst = rand.IntN(2) == 0
	}

	for player.HP > 0 && rival.HP > 0 {
		if playerActsFirst {
			move, err := g.choosePlayerMove(player, rival)
			if err != nil {
				return false, err
			}
			g.performMove(player, rival, move)
			g.renderBattle(player, rival, waitMedium)
			if player.HP <= 0 {
				return false, nil
			}
			if rival.HP <= 0 {
				return true, nil
			}

			g.performRivalTurn(rival, player)
			g.renderBattle(player, rival, waitMedium)
			if player.HP <= 0 {
				return false, nil
			}
			if rival.HP <= 0 {
				return true, nil
			}
		} else {
			g.performRivalTurn(rival, player)
			g.renderBattle(player, rival, waitMedium)
			if player.HP <= 0 {
				return false, nil
			}
			if rival.HP <= 0 {
				return true, nil
			}

			move, err := g.choosePlayerMove(player, rival)
			if err != nil {
				return false, err
			}
			g.performMove(player, rival, move)
			g.renderBattle(player, rival, waitMedium)
			if player.HP <= 0 {
				return false, nil
			}
			if rival.HP <= 0 {
				return true, nil
			}
		}
	}

	return player.HP > 0, nil
}

func (g *game) choosePlayerMove(player, rival *Monster) (Move, error) {
	if !hasUsableMove(player) {
		g.showMessage("Out of power!", "Kernel Panic activated!", waitMedium)
		return KernelPanic, nil
	}

	for {
		g.renderBattle(player, rival, waitNone)
		g.writePrompt("1, 2, 3, 4, Q")
		key, err := g.readKey()
		if err != nil {
			return Move{}, err
		}

		if key == 'q' {
			return Move{}, errQuit
		}
		if key < '1' || key > '4' {
			continue
		}

		index := int(key - '1')
		move := player.Moves[index]
		if move.Name == "" {
			continue
		}
		if move.PP <= 0 {
			g.showMessage("No power left!", "Choose another move.", waitMedium)
			continue
		}

		player.Moves[index].PP--
		return move, nil
	}
}

func (g *game) performRivalTurn(rival, player *Monster) {
	move := chooseComputerMove(rival, player)
	g.performMove(rival, player, move)
}

func chooseComputerMove(attacker, defender *Monster) Move {
	bestIndex := -1
	bestScore := -1.0

	for index, move := range attacker.Moves {
		if move.Name == "" || move.PP <= 0 {
			continue
		}

		effectiveness := getEffectiveness(move.Type, defender)
		score := float64(move.Power) * float64(move.Accuracy) / 100.0 * effectiveness
		if monsterHasType(attacker, move.Type) {
			score *= 1.25
		}
		score += float64(rand.IntN(16))

		if score > bestScore {
			bestScore = score
			bestIndex = index
		}
	}

	if bestIndex < 0 {
		return KernelPanic
	}

	attacker.Moves[bestIndex].PP--
	return attacker.Moves[bestIndex]
}

func (g *game) performMove(attacker, defender *Monster, move Move) {
	if attacker.HP <= 0 || defender.HP <= 0 {
		return
	}

	attackLine := attacker.Name + " used " + move.Name + "!"
	g.showMessage(attackLine, "", waitShort)

	if rand.IntN(100) >= move.Accuracy {
		g.showMessage(attacker.Name+"'s attack missed!", "", waitMedium)
		return
	}

	effectiveness := getEffectiveness(move.Type, defender)
	if effectiveness == 0 {
		g.showMessage("It had no effect!", "", waitMedium)
		return
	}

	damage := calculateDamage(attacker, defender, move)
	defender.HP = max(0, defender.HP-damage)

	if move.Recoil && damage > 0 {
		recoil := max(1, damage/4)
		attacker.HP = max(0, attacker.HP-recoil)
	}

	switch {
	case effectiveness > 1:
		g.showMessage("Super effective!", fmt.Sprintf("%s took %d damage!", defender.Name, damage), waitMedium)
	case effectiveness < 1:
		g.showMessage("Not very effective...", fmt.Sprintf("%s took %d damage!", defender.Name, damage), waitMedium)
	default:
		g.showMessage(fmt.Sprintf("%s took %d damage!", defender.Name, damage), "", waitMedium)
	}
}

func calculateDamage(attacker, defender *Monster, move Move) int {
	if attacker == nil || defender == nil || attacker.HP <= 0 || defender.HP <= 0 || move.Power <= 0 {
		return 0
	}

	effectiveness := getEffectiveness(move.Type, defender)
	if effectiveness == 0 {
		return 0
	}

	attack := attacker.Stats.Attack
	defense := defender.Stats.Defense
	if move.Class == DamageSpecial {
		attack = attacker.Stats.SpecialAttack
		defense = defender.Stats.SpecialDefense
	}
	if defense < 1 {
		defense = 1
	}

	levelFactor := float64((2*attacker.Level)/5 + 2)
	stab := 1.0
	if monsterHasType(attacker, move.Type) {
		stab = 1.25
	}
	randomFactor := float64(rand.IntN(16)+85) / 100.0
	damage := (((levelFactor * float64(move.Power) * float64(attack) / float64(defense)) / 50.0) + 2.0) * stab * effectiveness * randomFactor
	return max(1, int(damage))
}

func getEffectiveness(moveType MonsterType, defender *Monster) float64 {
	if defender == nil {
		return 1
	}

	effectiveness := 1.0
	for _, defenderType := range defender.Types {
		if defenderType == "" {
			continue
		}
		if modifier, ok := typeChart[moveType][defenderType]; ok {
			effectiveness *= modifier
		}
	}
	return effectiveness
}

func hasUsableMove(monster *Monster) bool {
	for _, move := range monster.Moves {
		if move.Name != "" && move.PP > 0 {
			return true
		}
	}
	return false
}

func monsterHasType(monster *Monster, target MonsterType) bool {
	if monster == nil {
		return false
	}
	for _, monsterType := range monster.Types {
		if monsterType == target {
			return true
		}
	}
	return false
}

func (g *game) writePrompt(options string) {
	optionLength := utf8.RuneCountInString(options)
	if optionLength > promptOptionWidth {
		options = truncateRunes(options, promptOptionWidth)
		optionLength = promptOptionWidth
	}

	padding := promptOptionWidth - optionLength
	leftPadding := padding / 2
	rightPadding := padding - leftPadding
	fmt.Fprintf(
		g.output,
		"%s(%s%s%s) →%s ",
		ansiDim,
		strings.Repeat(" ", leftPadding),
		options,
		strings.Repeat(" ", rightPadding),
		ansiReset,
	)
}

func formatNameLevel(name string, level int) string {
	levelLabel := fmt.Sprintf("Lv.%d", level)
	name = truncateRunes(name, screenWidth-utf8.RuneCountInString(levelLabel)-2)
	spacing := screenWidth - utf8.RuneCountInString(name) - utf8.RuneCountInString(levelLabel)
	return name + strings.Repeat(" ", spacing) + levelLabel
}

func formatHealthLine(currentHP, maxHP int) string {
	bar := getHealthBar(currentHP, maxHP)
	healthValue := max(currentHP, 0)
	if maxHP > 0 {
		healthValue = min(healthValue, maxHP)
	} else {
		healthValue = 0
	}
	value := fmt.Sprintf("%d/%d", healthValue, max(maxHP, 0))
	value = truncateRunes(value, screenWidth-utf8.RuneCountInString(bar))
	spacing := screenWidth - utf8.RuneCountInString(bar) - utf8.RuneCountInString(value)
	return bar + strings.Repeat(" ", spacing) + value
}

func colorControllerControls(line string) string {
	runes := []rune(line)
	var builder strings.Builder
	builder.Grow(len(line) + 64)
	activeColor := ""

	for index, character := range runes {
		if character != '█' {
			if activeColor != "" {
				builder.WriteString(ansiReset)
				activeColor = ""
			}
			builder.WriteRune(character)
			continue
		}

		color := ansiDPadBlack
		if index >= consoleButtonStart {
			color = ansiButtonPink
		}
		if color != activeColor {
			if activeColor != "" {
				builder.WriteString(ansiReset)
			}
			builder.WriteString(color)
			activeColor = color
		}
		builder.WriteRune(character)
	}

	if activeColor != "" {
		builder.WriteString(ansiReset)
	}
	return builder.String()
}

func (g *game) renderBattle(player, rival *Monster, wait waitDuration) {
	screen := blankScreen()
	screen[1] = formatNameLevel(rival.Name, rival.Level)
	screen[2] = formatHealthLine(rival.HP, rival.Stats.HP)
	screen[3] = strings.Repeat("─", screenWidth)
	screen[4] = formatNameLevel(player.Name, player.Level)
	screen[5] = formatHealthLine(player.HP, player.Stats.HP)
	screen[6] = formatMoveRow(player.Moves[0], 1, player.Moves[1], 2)
	screen[7] = formatMoveRow(player.Moves[2], 3, player.Moves[3], 4)
	g.renderGameBoy(screen, wait)
}

func formatMoveRow(first Move, firstNumber int, second Move, secondNumber int) string {
	return fitCell(formatMoveChoice(first, firstNumber), screenWidth/2) + fitCell(formatMoveChoice(second, secondNumber), screenWidth/2)
}

func formatMoveChoice(move Move, number int) string {
	if move.Name == "" {
		return ""
	}
	return fmt.Sprintf("[%d] %s", number, move.Name)
}

func (g *game) showMessage(first, second string, wait waitDuration) {
	screen := blankScreen()
	screen[3] = centerLine(first)
	screen[4] = centerLine(second)
	g.renderGameBoy(screen, wait)
}

func (g *game) showVictory(winner Monster) {
	screen := blankScreen()
	screen[1] = centerLine("Hall of Fame")
	screen[3] = centerLine(winner.Name + "  Lv." + strconv.Itoa(winner.Level))
	screen[5] = centerLine("Congratulations!")
	screen[6] = centerLine("New champion!")
	g.renderGameBoy(screen, waitLong)
}

func (g *game) showDefeat(loser Monster) {
	screen := blankScreen()
	screen[3] = centerLine(loser.Name + " fainted!")
	screen[4] = centerLine("The terminal goes quiet...")
	g.renderGameBoy(screen, waitLong)
}

func (g *game) readKey() (rune, error) {
	fmt.Fprint(g.output, ansiHidden)
	line, err := g.input.ReadString('\n')
	fmt.Fprint(g.output, ansiReset)
	if err != nil && len(line) == 0 {
		return 0, err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return ' ', nil
	}

	key, _ := utf8.DecodeRuneInString(line)
	if key >= 'A' && key <= 'Z' {
		key += 'a' - 'A'
	}
	return key, nil
}

func (g *game) renderGameBoy(screen [screenHeight]string, wait waitDuration) {
	g.clearScreen()
	fmt.Fprintln(g.output)
	fmt.Fprintln(g.output, "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓")
	fmt.Fprintln(g.output, "▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓")
	fmt.Fprintln(g.output, "▓▓░░████████████████████████████████████████████░░▓▓")
	fmt.Fprintln(g.output, "▓▓░░██▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██░░▓▓")
	for _, line := range screen {
		fmt.Fprintf(g.output, "▓▓░░██▒▒ %s ▒▒██░░▓▓\n", fitLine(line))
	}
	fmt.Fprintln(g.output, "▓▓░░██▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██░░▓▓")
	fmt.Fprintln(g.output, "▓▓░░████████████████████████████████████████████░░▓▓")
	fmt.Fprintln(g.output, "▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓")
	fmt.Fprintln(g.output, "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓")
	controllerLines := [...]string{
		"▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒▒▒▒▒ ██ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒▒▒ ██████ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ████ ▒▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒████████████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████████▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒▒▒ ██████ ▒▒▒▒▒▒▒▒▒▒▒ ████ ▒▒▒▒▒ ████ ▒▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒▒▒▒▒ ██ ▒▒▒▒▒▒▒▒▒▒▒▒████████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ████ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒ ██████ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓",
		"▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒ ██████ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓",
		"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓",
	}
	for _, line := range controllerLines {
		fmt.Fprintln(g.output, colorControllerControls(line))
	}

	if wait > 0 {
		time.Sleep(time.Duration(wait))
	}
}

func (g *game) clearScreen() {
	fmt.Fprint(g.output, "\033[H\033[2J")
}

func (g *game) resetTerminal() {
	fmt.Fprintf(g.output, "%s\033[?25h\n", ansiReset)
}

func titleScreen() [screenHeight]string {
	screen := blankScreen()
	screen[1] = centerLine("B A S H M O N")
	screen[3] = centerLine("Glitch Version")
	screen[5] = centerLine("© Mockware 2026")
	return screen
}

func blankScreen() [screenHeight]string {
	var screen [screenHeight]string
	screen[0] = strings.Repeat("═", screenWidth)
	screen[screenHeight-1] = strings.Repeat("═", screenWidth)
	for index := 1; index < screenHeight-1; index++ {
		screen[index] = strings.Repeat(" ", screenWidth)
	}
	return screen
}

func getHealthBar(currentHP, maxHP int) string {
	if maxHP <= 0 || currentHP <= 0 {
		return strings.Repeat("░", 10)
	}
	if currentHP > maxHP {
		currentHP = maxHP
	}

	filled := currentHP * 10 / maxHP
	if filled == 0 {
		filled = 1
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", 10-filled)
}

func formatMonsterTypes(monster Monster) string {
	if monster.Types[1] == "" {
		return string(monster.Types[0])
	}
	return string(monster.Types[0]) + " / " + string(monster.Types[1])
}

func centerLine(content string) string {
	content = truncateRunes(content, screenWidth)
	length := utf8.RuneCountInString(content)
	left := (screenWidth - length) / 2
	right := screenWidth - length - left
	return strings.Repeat(" ", left) + content + strings.Repeat(" ", right)
}

func fitLine(content string) string {
	content = truncateRunes(content, screenWidth)
	length := utf8.RuneCountInString(content)
	return content + strings.Repeat(" ", screenWidth-length)
}

func fitCell(content string, width int) string {
	content = truncateRunes(content, width)
	length := utf8.RuneCountInString(content)
	return content + strings.Repeat(" ", width-length)
}

func truncateRunes(content string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(content)
	if len(runes) > width {
		return string(runes[:width])
	}
	return content
}
