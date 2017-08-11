package output
import "encoding/binary"

const (
	aiffheaderID = aiffID{'F', 'O', 'R', 'M'}
	aiffheaderType = aiffID{'A', 'I', 'F', 'F'}
	aiffSoundChunkID = aiffID{'S', 'S', 'N', 'D'}
	aiffCommonChunkID = aiffID{'C', 'O', 'M', 'M'}
)

type AIFF struct {
	chunks []aiffChunk
}


func NewAIFFOutput(name, author, copywrite string, notes []string) *AIFF {
	aiff := new(AIFF)
	aiff.chunks = append(aiff.chunks, newAIFFCommonChunk())
	if len(name) > 0 {
		aiff.chunks = append(aiff.chunks, newAIFFNameChunk(name))
	}
	if len(author) > 0 {
		aiff.chunks = append(aiff.chunks, newAIFFAuthorChunk(name))
	}
	if len(copywrite) > 0 {
		aiff.chunks = append(aiff.chunks, newAIFFCopywriteChunk(name))
	}
	for _, note := range(notes) {
		aiff.chunks = append(aiff.chunks, newAIFFAnnotationChunk(note))
	}
	aiff.chunks = append(aiff.chunks, newAIFFAnnotationChunk(madeWithSynthiaString))
	return aiff
}


func (a *AIFF) WritePCM(p []int32, samplerate float64) error {
	
}

type aiffID [4]byte
type aiffShort int16
type aiffLong int32
type aiffUnsignedShort uint16
type aiffUnsignedLong uint32
type aiffFloat [10]byte

type aiffChunk interface {
	id() aiffID
	chunkSize() aiffLong
	write()
}

type aiffCommonChunk struct {
	size aiffLong
	numChannels aiffShort
	numFrames	aiffUnsignedLong
	sampleSize	aiffShort
	samplerate  aiffFloat	
}

func (c *aiffCommonChunk) id() aiffID{
	return aiffCommonChunkID
}

func (c *aiffCommonChunk) chunkSize() aiffLong {
	return c.size
}

func newAIFFCommonChunk() *aiffCommonChunk {
	c := new(aiffCommonChunk)
	c.sampleSize = pcmSampleSize
	return c
}

type aiffSoundChunk struct {
	size aiffLong
	offset aiffUnsignedLong
	blockSize aiffUnsignedLong
	data []byte
}

func (c *aiffSoundChunk) id() aiffID{
	return aiffSoundChunkID
}

func (c *aiffSoundChunk) chunkSize() aiffLong {
	return c.size
}

type aiffStringChunk struct {
	id aiddID
	data string
}

func (c *aiffStringChunk) id() aiffID{
	return c.id()
}

func (c *aiffStringChunk) chunkSize() aiffLong {
	return aiffLong(len(c.data))
}

func newAIFFNameChunk(name string) *aiffStringChunk {
	c := new(aiffStringChunk)
	c.id = aiffID{'N', 'A', 'M', 'E'}
	c.data = name
	return aiffStringChunk
}

func newAIFFCopywriteChunk(copywrite string) *aiffStringChunk {
	c := new(aiffStringChunk)
	c.id = aiffID{'(', 'c', ')', ' '}
	c.data = copywrite
	return aiffStringChunk
}

func newAIFFAuthorChunk(author string) *aiffStringChunk {
	c := new(aiffStringChunk)
	c.id = aiffID{'A', 'U', 'T', 'H'}
	c.data = author
	return aiffStringChunk
}

func newAIFFAnnotationChunk(annotation string) *aiffStringChunk {
	c := new(aiffStringChunk)
	c.id = aiffID{'A', 'N', 'N', 'O'}
	c.data = annotation
	return aiffStringChunk
}

func float2AIFF(f float64) aiffFloat {
	
}

func aiff2Float(af aiffFloat) float64 {
		
}