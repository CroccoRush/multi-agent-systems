package experiment

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"

	"multi-agent-systems/t1/internal/agent"
	env "multi-agent-systems/t1/internal/environment"
)

// Topology представляет топологию агентов
type Topology struct {
	Agents []struct {
		ID        int      `yaml:"id"`
		Number    *float64 `yaml:"number,omitempty"`
		Neighbors []int    `yaml:"neighbors"`
	} `yaml:"agents"`
}

// UnmarshalFile создаёт топологию из YAML-файла
func (t *Topology) UnmarshalFile(path string) error {
	topologyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading topology file: %s", err)
	}

	err = yaml.Unmarshal(topologyData, &t)
	if err != nil {
		return fmt.Errorf("error unmarshalling topology: %s", err)
	}

	return nil
}

func Run() {
	var topology Topology
	if err := topology.UnmarshalFile("topology.yaml"); err != nil {
		fmt.Println("Error unmarshalling topology:", err)
		return
	}

	environments := env.NewEnvironments(
		len(topology.Agents),
		0,
		make(map[int][]env.AgentLink),
	)

	// Создаём агентов
	totalNumber := float64(0)
	agents := make(map[int]*agent.Agent)
	for _, agentData := range topology.Agents {
		agents[agentData.ID] = agent.NewAgent(agentData.ID, agentData.Number)
		totalNumber += agents[agentData.ID].OwnValue()
	}

	// Создаём агентов друг с другом по топологии
	for _, agentData := range topology.Agents {
		links := make([]env.AgentLink, len(agentData.Neighbors))
		for i, neighborID := range agentData.Neighbors {
			links[i] = env.AgentLink{
				ID:      neighborID,
				Channel: agents[neighborID].Channel(),
			}
		}
		environments.AddAgentLinks(agentData.ID, links)
	}
	environments.SetAverage(totalNumber / float64(environments.AgentCount()))

	fmt.Printf(
		"Agent count = %d --- Average = %f\n",
		environments.AgentCount(),
		environments.Average(),
	)

	// Запускаем агентов
	var wg sync.WaitGroup
	wg.Add(len(agents))
	for _, a := range agents {
		go a.Run(&wg, &environments)
	}
	// Ждем завершения работы всех агентов
	wg.Wait()
}
