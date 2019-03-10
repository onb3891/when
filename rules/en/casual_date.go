package en

import (
	"regexp"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/olebedev/when/rules"
)

func CasualDate(s rules.Strategy) rules.Rule {
	overwrite := s == rules.Override

	return &rules.F{
		RegExp: regexp.MustCompile("(?i)(?:\\W|^)(now|today|tonight|last\\s*night|last\\s*day|last\\s*year|last\\s*month|last\\s*week|this\\s*year|(?:tomorrow|tmr|yesterday)\\s*|tomorrow|tmr|yesterday)(?:\\W|$)"),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			lower := strings.ToLower(strings.TrimSpace(m.String()))

			switch {
			case strings.Contains(lower, "tonight"):
				if c.Hour == nil && c.Minute == nil || overwrite {
					c.Hour = pointer.ToInt(23)
					c.Minute = pointer.ToInt(0)
				}
			case strings.Contains(lower, "today"):
				// c.Hour = pointer.ToInt(18)
			case strings.Contains(lower, "tomorrow"), strings.Contains(lower, "tmr"):
				if c.Duration == 0 || overwrite {
					c.Duration += time.Hour * 24
				}
			case strings.Contains(lower, "yesterday"):
				if c.Duration == 0 || overwrite {
					c.Duration -= time.Hour * 24
				}
			case strings.Contains(lower, "last night"):
				if (c.Hour == nil && c.Duration == 0) || overwrite {
					c.Hour = pointer.ToInt(23)
					c.Duration -= time.Hour * 24
				}
			case strings.Contains(lower, "last day"):
				if (c.Hour == nil && c.Duration == 0) || overwrite {
					c.Hour = pointer.ToInt(00)
					c.Duration -= time.Hour * 24
				}
			case strings.Contains(lower, "last week"):
				if (c.Day == nil && c.Duration == 0) || overwrite {
					t := ref.AddDate(0, 0, -(int(ref.Weekday())))
					t = t.AddDate(0, 0, -7)
					c.Day = pointer.ToInt((t.Day()))
					c.Year = pointer.ToInt(t.Year())
					c.Month = pointer.ToInt(int(t.Month()))
					c.Hour = pointer.ToInt(00)
					//c.Duration += time.Hour * 24 * 7
				}
			case strings.Contains(lower, "last year"):
				if (c.Hour == nil && c.Duration == 0) || overwrite {
					c.Year = pointer.ToInt(ref.Year() - 1)
					c.Month = pointer.ToInt(1)
					c.Day = pointer.ToInt(1)
					c.Duration -= time.Hour * 24
				}
			case strings.Contains(lower, "last month"):
				if (c.Hour == nil && c.Duration == 0) || overwrite {
					t := ref.AddDate(0, -1, 0)

					c.Day = pointer.ToInt((t.Day()))
					c.Year = pointer.ToInt(t.Year())
					c.Month = pointer.ToInt(int(t.Month()))
					c.Hour = pointer.ToInt(00)
					c.Minute = pointer.ToInt(01)
					//c.Duration -= time.Hour * 24
				}
			case strings.Contains(lower, "this year"):
				if (c.Hour == nil && c.Duration == 0) || overwrite {
					c.Day = pointer.ToInt(1)
					c.Year = pointer.ToInt(ref.Year())
					c.Month = pointer.ToInt(1)
					c.Hour = pointer.ToInt(00)
					c.Minute = pointer.ToInt(01)
					//c.Duration -= time.Hour * 24
				}
			}

			return true, nil
		},
	}
}
