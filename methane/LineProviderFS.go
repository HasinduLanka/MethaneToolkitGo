package methane

import "bufio"

func NewReusableLineProvider_FromURI(uri string) (*ReusableLineProvider, error) {
	filename := NewCacheFilename()
	dwnErr := DownloadToFile(filename, uri)

	if PrintError(dwnErr) {
		return nil, dwnErr
	}

	return NewReusableLineProvider_FromDiscardableFile(filename)
}

func NewLineProvider_FromURI(uri string) (*LineProvider, error) {
	F, err := LoadURIToString(uri)

	if err != nil {
		return nil, err
	}

	return NewLineProvider_FromString(F), nil
}

func NewLineProvider_FromFile(filename string) (*LineProvider, error) {
	file, err := LoadFileToIOReader(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	LP := &LineProvider{make(chan string, 256)}

	go func() {
		for scanner.Scan() {
			LP.ch <- scanner.Text()
		}

		close(LP.ch)
		_ = file.Close()
	}()

	return LP, nil
}

func NewReusableLineProvider_FromURI_DiskCached(uri string) (*ReusableLineProvider, error) {
	_, err := LoadURIToStringCached(uri)

	if err != nil {
		return nil, err
	}
	filename := CachedName(uri)

	return NewReusableLineProvider_FromDiscardableFile(filename)
}

func NewReusableLineProvider_FromDiscardableFile(filename string) (*ReusableLineProvider, error) {
	LPr, err := NewReusableLineProvider_FromFile(filename)

	if err != nil {
		LPr.Discard = func() {
			DeleteFiles(filename)
		}
	}
	return LPr, err
}

func NewReusableLineProvider_FromFile(filename string) (*ReusableLineProvider, error) {
	return &ReusableLineProvider{
		GetInstance: func() *LineProvider {
			P, PErr := NewLineProvider_FromFile(filename)
			if PErr != nil {
				return NewLineProvider_FromString("")
			}
			return P
		},
	}, nil
}

func (lp *LineProvider) CacheToDisk() (*ReusableLineProvider, error) {
	filename := NewCacheFilename()

	for line := range lp.ch {
		AppendFile(filename, []byte(line+"\n"))
	}

	return NewReusableLineProvider_FromDiscardableFile(filename)
}

func NewByteSource_FromURI(uri string) *ByteSource {
	BS := ByteSource{
		Get: func() []byte {
			B, err := LoadURI(uri)
			if err != nil {
				Print("ByteSourceFromURI Error : " + uri)
				PrintError(err)
				return []byte{}
			}
			return B
		},
	}
	return &BS
}

func URIsToByteSources(uris *LineProvider) chan *ByteSource {
	var sources chan *ByteSource = make(chan *ByteSource)
	go func() {
		for line := range uris.ch {
			sources <- NewByteSource_FromURI(line)
		}
		close(sources)
	}()
	return sources
}
