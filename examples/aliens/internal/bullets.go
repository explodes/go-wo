package internal

import (
	"sync"

	"github.com/explodes/go-wo/wobj"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Bullet struct {
	*wobj.Object
}

var (
	bulletBounds = pixel.R(0, 0, 6, 6)
)

type BulletPool struct {
	Active map[*Bullet]struct{}
	pool   *sync.Pool
}

func NewBulletPool() *BulletPool {
	drawable := newBulletDrawable()
	factory := func() interface{} {
		return &Bullet{
			Object: wobj.NewDrawableObject(drawable, 0, 0, 6, 6),
		}
	}
	return &BulletPool{
		pool:   &sync.Pool{New: factory},
		Active: make(map[*Bullet]struct{}),
	}
}

func (b *BulletPool) Spawn() *Bullet {
	i := b.pool.Get()
	bullet := i.(*Bullet)
	b.Active[bullet] = struct{}{}
	return bullet
}

func (b *BulletPool) Return(bullet *Bullet) {
	b.pool.Put(bullet)
	delete(b.Active, bullet)
}

type bulletDrawable struct {
	im *imdraw.IMDraw
}

func newBulletDrawable() *bulletDrawable {
	im := imdraw.New(nil)

	im.Color = colornames.Red
	im.Push(pixel.V(3, 3))
	im.Circle(3, 0)

	return &bulletDrawable{im: im}
}

func (b *bulletDrawable) Draw(target pixel.Target, mat pixel.Matrix) {
	c := target.(*pixelgl.Canvas)
	c.SetMatrix(mat)
	b.im.Draw(target)
	c.SetMatrix(pixel.IM)
}

func (b *bulletDrawable) Bounds() pixel.Rect {
	return bulletBounds
}
