package audio

import (
	"fmt"
	"io"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Player interface {
	Play(audio io.Reader)
	Close()
}

const maxQueueSize = 20

type audioPlayer struct {
	ctx   *oto.Context
	queue chan io.Reader
	done  chan struct{}
}

func NewAudioPlayer() (Player, error) {
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}
	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-ready
	player := &audioPlayer{
		ctx:   ctx,
		queue: make(chan io.Reader, maxQueueSize),
		done:  make(chan struct{}),
	}
	go player.worker()
	return player, nil
}

func (p *audioPlayer) worker() {
	for {
		select {
		case file := <-p.queue:
			p.play(file)
		case <-p.done:
			close(p.queue)
			return
		}
	}
}

func (p *audioPlayer) play(audio io.Reader) {
	decodedMp3, err := mp3.NewDecoder(audio)
	if err != nil {
		fmt.Println("Error decoding audio:", err)
		return
	}

	stream := p.ctx.NewPlayer(decodedMp3)
	stream.Play()
	for stream.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}
}

func (p *audioPlayer) Play(audio io.Reader) {
	p.queue <- audio
}

func (p *audioPlayer) Close() {
	close(p.done)
}
