package presenters

const (
	MarkerStringCropped = "..." // Маркер-признак того, что строка "обрезана".
)

var _ StringPresenterWriter = (*Cropper)(nil)

// Cropper позволяет записывать в него информацию, но после заданного порога запись в буфер прекращается.
// Дальнейшая запись игнорируется. Если пороговое значение было превышено, то при получении строкового представления
// в конце строки появится многоточие.
// Например, полезно при записи логов обрезать строку до указанного значения. Чтобы не перегружать систему выводом
// супер-больших объектов в лог.
// Внимание! Длина вычисляется не в символах, а в байтах!
type Cropper struct {
	buf      []byte
	limit    int
	exceeded bool
}

func (c *Cropper) Write(chunk []byte) (int, error) {
	want := len(chunk)
	can := want
	left := c.limit - len(c.buf)

	if want > left {
		can = left
		c.exceeded = true
	}

	if can > 0 {
		c.buf = append(c.buf, chunk[:can]...)
		if c.exceeded {
			c.buf = append(c.buf, []byte(MarkerStringCropped)...)
		}
	}

	return want, nil
}

func (c *Cropper) String() string {
	return string(c.buf)
}

// NewCropper создает новый экземпляр Cropper с установленным пороговым значением.
func NewCropper(maxBytes int) *Cropper {
	return &Cropper{
		limit: maxBytes,
		buf:   make([]byte, 0, maxBytes+len(MarkerStringCropped)),
	}
}

// StringTail вернет "хвост" строки не длиннее заданного размера.
func StringTail(value string, tailLength int) string {
	if value == "" {
		return ""
	}

	l := len(value) - tailLength
	if l < 0 {
		l = 0
	}

	return value[l:]
}
