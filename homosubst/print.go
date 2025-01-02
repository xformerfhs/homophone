package homosubst

import "fmt"

func (s *Substitutor) Print() {
	substitutions := s.substitutions
	for i, l := range substitutions {
		fmt.Printf(`   %c: `, i+'A')
		for _, r := range l {
			fmt.Printf(`%c`, r)
		}
		fmt.Println()
	}
}
