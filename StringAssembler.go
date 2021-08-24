package main

type StringAssembler struct {
	Providers []ReusableLineProvider
	maxDepth  int
}

func NewStringAssembler() StringAssembler {
	return StringAssembler{
		Providers: []ReusableLineProvider{},
	}
}

func (s *StringAssembler) Add(provider *ReusableLineProvider) {
	s.Providers = append(s.Providers, *provider)
}
func (s *StringAssembler) AddStatic(strng string) {
	s.Providers = append(s.Providers, *NewReusableLineProvider_Static(strng))
}

func (s *StringAssembler) Assemble() *ReusableLineProvider {
	return s.AssembleThreaded(0)
}
func (s *StringAssembler) AssembleThreaded(threadDepth int) *ReusableLineProvider {
	s.maxDepth = len(s.Providers)

	return &ReusableLineProvider{
		GetInstance: func() *LineProvider {
			LP := &LineProvider{
				ch: make(chan string, 64),
			}

			chcloser := make(chan bool, 64)
			chcloser <- true
			go s._assembleRecursive(0, "", &LP.ch, &chcloser, threadDepth)

			go func() {
				count := 0
				for cls := range chcloser {
					if cls {
						count++
					} else {
						count--

						if count <= 0 {
							close(LP.ch)
						}
					}
				}
			}()

			return LP
		},
	}
}

func (s *StringAssembler) _assembleRecursive(iProv int, build string, ch *chan string, chcloser *chan bool, threadDepth int) {

	if iProv >= s.maxDepth {
		*ch <- build
		*chcloser <- false
		return
	}

	provider := s.Providers[iProv]
	iProv += 1

	if threadDepth > 0 {
		threadDepth--

		for line := range provider.GetInstance().ch {
			cbuild := build + line

			*chcloser <- true
			go s._assembleRecursive(iProv, cbuild, ch, chcloser, threadDepth)
		}
	} else {
		for line := range provider.GetInstance().ch {
			cbuild := build + line

			*chcloser <- true
			s._assembleRecursive(iProv, cbuild, ch, chcloser, threadDepth)
		}
	}

	*chcloser <- false

}
