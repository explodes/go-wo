package wo

import (
	"bytes"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/pkg/errors"
)

const (
	audioSampleRate      = beep.SampleRate(44100)
	audioResampleQuality = 4
	audioBufferTiming    = time.Second / 10
)

// Sound is a slice of bytes that can be decoded
// and played on a Speaker.
type Sound struct {
	samples     []byte
	format      string
	audioFormat beep.Format
}

// NewSound creates a new sound from some bytes. It is decoded
// either as "wav" or "mp3" as specified.
func NewSound(format string, samples []byte) (*Sound, error) {
	_, audioFormat, err := decode(format, samples)
	if err != nil {
		return nil, err
	}
	sound := &Sound{
		samples:     samples,
		format:      format,
		audioFormat: audioFormat,
	}
	return sound, nil
}

// stream creates a new stream from a Sound so that it
// can be played on a Speaker.
func (s *Sound) stream() beep.Streamer {
	stream, _, err := decode(s.format, s.samples)
	if err != nil {
		// it decoded once before, how did it not decode correctly a second time?
		panic(errors.Errorf("secondary decoding failed: %v", err))
	}
	return stream
}

// decode attempts to decode a sound in the given format.
func decode(format string, samples []byte) (beep.Streamer, beep.Format, error) {
	reader := &readCloserWrapper{Reader: bytes.NewReader(samples)}
	switch format {
	case "wav":
		return wav.Decode(reader)
	case "mp3":
		return mp3.Decode(reader)
	default:
		return nil, beep.Format{}, errors.Errorf("audio format %s not supported", format)
	}
}

// Audible is a sound already connected to a Speaker
// for convenience so that calling Play() will play
// the sound.
type Audible struct {
	speaker *Speaker
	sound   *Sound
}

// Play plays this Audible's Sound on its Speaker.
func (a *Audible) Play() {
	a.speaker.Play(a.sound)
}

// Speaker is the output for a Sound.
type Speaker struct {
	//sampleRate   beep.SampleRate
	//quality      int
	//bufferTiming time.Duration

	mixer *beep.Mixer
}

// NewSpeaker creates and initializes a new Speaker.
func NewSpeaker() (*Speaker, error) {
	spkr := &Speaker{
		mixer: &beep.Mixer{},
	}
	err := spkr.init()
	if err != nil {
		return nil, err
	}
	return spkr, nil
}

// init prepares this Speaker for playback.
func (s *Speaker) init() error {
	err := speaker.Init(audioSampleRate, audioSampleRate.N(audioBufferTiming))
	if err != nil {
		return err
	}
	speaker.Play(s.mixer)
	return nil
}

// Play plays a sound on this speaker.
func (s *Speaker) Play(sound *Sound) {
	stream := s.ensureSampleRate(sound)
	s.mixer.Play(stream)
}

// ensureSampleRate makes sure that the sound is already formatted
// in the same format as this speaker. If it is not, it is resampled
// as it plays.
func (s *Speaker) ensureSampleRate(sound *Sound) beep.Streamer {
	stream := sound.stream()
	if sound.audioFormat.SampleRate != audioSampleRate {
		return beep.Resample(audioResampleQuality, sound.audioFormat.SampleRate, audioSampleRate, stream)
	}
	return stream
}

// Audible creates a new Audible from a Sound for this Speaker.
func (s *Speaker) Audible(sound *Sound) *Audible {
	return &Audible{
		speaker: s,
		sound:   sound,
	}
}
