package notification

type Notify struct {

    Title       string
    Message     string
    Icon        string
    Urgency     string
    ExpireTime  int
    Category    string
    Hint       string
}

const (

    LowPriority      = "low"
    NormalPriority   = "normal"
    CriticalPriority = "critical"
)