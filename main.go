package main

import (
    "encoding/xml"
    "fmt"
    "os"
    "log"
    "os/exec"
    "strings"
    "strconv"
    "bufio"
    "regexp"
   )

func alpha_num(str string) string {
    return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "")
}

func to_int(str string) int {
    integer, _ := strconv.Atoi(str)
    return integer
}

func up_or_down(num int, direction int) int {
    if direction % 2 == 0 {
        return num
    } else {
        return -num
    }
}

func pitch_and_octave(num int) []string {
    pitches := []string{"C_","C#","D_","D#","E_","F_","F#","G_","G#","A_","A#","B_"}
    for i,_ := range pitches {
        pitches[i] = pitches[i] + fmt.Sprint(num)

    }
    return pitches
}

func pitch_list() []string {
    notes := []string{}
    for i:=0; i <=7; i++ {
        notes = append(notes, pitch_and_octave(i)...)
    }
    return notes
}

func indexOf(element string, data []string) int {
    for k, v := range data {
        if element == v {
            return k
        }
    }
    return -1
}

func accidental(accidental string) string {
        if accidental == "#" {
            return "1"
        } else if accidental == "b" {
            return "-1"
        } 
        return "0"
    }

func pitchXML(note string) []Pitch {
    split := strings.Split(note,"")
    return []Pitch{
        {
            Step: split[0],
            Alter: accidental(split[1]),
            Octave: split[2],
        },
    }
}


func generate_left_note(note string, interval int, duration int, notetype string, notes []string) interface{} {
    var result interface{}
    if strings.Contains(note,"rest"){
        result = rh(note, duration, notetype)
    } else if ! strings.Contains(note,"rest") && indexOf(note,notes)+interval < len(notes)   {
        result = rhc(note ,interval,duration,notetype,notes)
    } else if ! strings.Contains(note,"rest") && indexOf(note,notes)+interval >= len(notes)  {
        result = rh(note, duration, notetype)
    }
    return result
}



func next_note(notes []string, original_note string, interval int, direction int, lower_bound int, upper_bound int, flag int) (string, int) {
    if strings.Contains(original_note, "rest-") {
        original_note = original_note[5:]
    }
    next_note_index := indexOf(original_note, notes)+up_or_down(interval, direction)

    if next_note_index > upper_bound && flag == 0 {
        return "rest-"+original_note, 1
    } else if next_note_index > upper_bound && flag == 1 {
        pitch := notes[next_note_index-12]
        return pitch, 0
    }
    if next_note_index < lower_bound && flag == 0 {
        return "rest-"+original_note, 1
    } else if next_note_index < lower_bound && flag == 1 {
        pitch := notes[next_note_index+12]
        return pitch, 0
    } else if next_note_index >= lower_bound && next_note_index <= upper_bound {
        pitch := notes[next_note_index]
        return pitch, 0
    }
    return "error", -1
}


func rhythm(num int) string {
    rhythm_name := map[int]string{
    1: "eighth",
    2: "quarter",
    3: "quarter",
    4: "half",
    5: "half", // with tied eighth
    6: "half", // with tied quarter
    7: "half", // with tied dotted quarter
    8: "whole", 
    9: "whole", // with tied eighth
    10: "whole", // with tied quarter
}
    return rhythm_name[num]
}



func note_append(note_array []Note, in interface{}) []Note {
    switch in.(type) {
    default: 
        return note_array
        
    case Note:

        return append(note_array, in.(Note))
    case []Note:
        for index, _ := range in.([]Note) {
            note_array = append(note_array, in.([]Note)[index])

        }
        return note_array
    }

}


func rh(note string, duration int, notetype string)  interface{} {
    if strings.Contains(note,"rest") {

        switch duration {
            case 1:
                return Note{
                    Rest: "HERE",
                    Duration: strconv.Itoa(1),
                    Type:     "eighth",
                    }
            case 2:
                return Note{
                    Rest: "HERE",
                    Duration: strconv.Itoa(2),
                    Type:     "quarter",
                    }

            case 3:
            return Note{
                Rest: "HERE",
                Duration: strconv.Itoa(3),
                Type:     "quarter",
                Dot: "HERE",
                }
        case 4:
            return Note{
                Rest: "HERE",
                Duration: strconv.Itoa(4),
                Type:     "half",
                }
        case 5:
            return []Note{
                {
                Rest: "HERE",
                Duration: strconv.Itoa(4),
                Type: "half",
                },
                {
                Rest: "HERE",
                Duration: strconv.Itoa(1),
                Type: "eighth",
                },

            }
        case 6:
            return Note{
                Rest: "HERE",
                Duration: strconv.Itoa(duration),
                Type:     rhythm(duration),
                Dot: "HERE",
                }
        case 7:
            return []Note{
                {
                Rest: "HERE",
                Duration: strconv.Itoa(4),
                Type: "half",
                },
                {
                Rest: "HERE",
                Duration: strconv.Itoa(3),
                Type: "quarter",
                Dot: "HERE",
                },

            }
        case 8:
            return Note{
                Rest: "HERE",
                Duration: strconv.Itoa(8),
                Type:     "half",
                }
        case 9:
            return []Note{
                {
                Rest: "HERE",
                Duration: strconv.Itoa(8),
                Type: "whole",
                },
                {
                Rest: "HERE",
                Duration: strconv.Itoa(1),
                Type: "eighth",
                },

            }
        case 10:
            return []Note{
                {
                Rest: "HERE",
                Duration: strconv.Itoa(8),
                Type: "whole",
                },
                {
                Rest: "HERE",
                Duration: strconv.Itoa(2),
                Type: "quarter",
                },

            }
        default:
            fmt.Println("error")
            return -1
        }
    } else {

        switch duration {
        case 1:
            return Note{
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(1),
                Type:     "eighth",
                }
        case 2:
            return Note{
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(2),
                Type:     "quarter",
                }

        case 3:
            return Note{
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(3),
                Type:     "quarter",
                Dot: "HERE",
                }
        case 4:
            return Note{
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(4),
                Type:     "half",
                }
        case 5:
            return []Note{
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(4),
                Tie: []Tie{
                    {
                        Type: "start",
                    },
                },
                Type: "half",
                },
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(1),
                Tie: []Tie{
                    {
                        Type: "stop",
                    },
                },
                Type: "eighth",
                },

            }
        case 6:
            return Note{
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(duration),
                Type:     rhythm(duration),
                Dot: "HERE",
                }
        case 7:
            return []Note{
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(4),
                Tie: []Tie{
                    {
                        Type: "start",
                    },
                },
                Type: "half",
                },
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(3),
                Tie: []Tie{
                    {
                        Type: "stop",
                    },
                },
                Type: "quarter",
                Dot: "HERE",
                },

            }
        case 8:
            return Note{
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(8),
                Type:     "whole",
                }
        case 9:
            return []Note{
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(8),
                Tie: []Tie{
                    {
                        Type: "start",
                    },
                },
                Type: "whole",
                },
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(1),
                Tie: []Tie{
                    {
                        Type: "stop",
                    },
                },
                Type: "eighth",
                },

            }
        case 10:
            return []Note{
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(8),
                Tie: []Tie{
                    {
                        Type: "start",
                    },
                },
                Type: "whole",
                },
                {
                Pitch: pitchXML(note),
                Duration: strconv.Itoa(2),
                Tie: []Tie{
                    {
                        Type: "stop",
                    },
                },
                Type: "quarter",
                },

            }
        default:
            fmt.Println("error")
            return -1
    }
    }
}


func rhc(note string, interval int, duration int, notetype string, notes []string)  interface{} {

        switch duration {
        case 1:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(1),
                    Type: "eighth",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(1),
                    Type: "eighth",
                },

            }
        case 2:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(2),
                    Type: "quarter",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(2),
                    Type: "quarter",
                },
            }

        case 3:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(3),
                    Type: "quarter",
                    Dot: "HERE",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(3),
                    Type: "quarter",
                    Dot: "HERE",
                },
            }
        case 4:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(4),
                    Type: "half",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(4),
                    Type: "half",
                },
            }
        case 5:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(4),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },
                    Type: "half",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(4),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },                    
                    Type: "half",
                },
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(1),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },
                    Type: "eighth",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(1),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },                    
                    Type: "eighth",
                },
            }
        case 6:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(6),
                    Type: "half",
                    Dot: "HERE",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(6),
                    Type: "half",
                    Dot: "HERE",
                },
            }
        case 7:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(4),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },
                    Type: "half",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(4),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },                    
                    Type: "half",
                },
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(3),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },
                    Type: "quarter",
                    Dot: "HERE",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(3),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },                    
                    Type: "quarter",
                    Dot: "HERE",
                },
            }
        case 8:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(8),
                    Type: "whole",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(8),
                    Type: "whole",
                },
            }
        case 9:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(8),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },
                    Type: "whole",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(8),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },                    
                    Type: "whole",
                },
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(1),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },
                    Type: "eighth",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(1),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },                    
                    Type: "eighth",
                },
            }
        case 10:
            return []Note{
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(8),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },
                    Type: "whole",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(8),
                    Tie: []Tie{
                        {
                            Type: "start",
                        },
                    },                    
                    Type: "whole",
                },
                {
                    Pitch: pitchXML(note),
                    Duration: strconv.Itoa(2),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },
                    Type: "quarter",
                },
                {
                    Chord: "HERE",
                    Pitch: pitchXML(notes[indexOf(note,notes)+interval]),
                    Duration: strconv.Itoa(2),
                    Tie: []Tie{
                        {
                            Type: "stop",
                        },
                    },                    
                    Type: "quarter",
                },
            }
        default:
            fmt.Println("error")
            return -1
    }
}

type ScorePartwise struct {
    XMLName xml.Name `xml:"score-partwise"`
    Version string `xml:"version,attr"`
    PartList []PartList `xml:"part-list"`
    Part     []Part `xml:"part"`
}

type PartList struct {
    XMLName xml.Name `xml:"part-list"`
    ScorePart []ScorePart `xml:"score-part"`
}

type ScorePart struct {
    XMLName xml.Name `xml:"score-part"`
    Id string `xml:"id,attr"`
    PartName string `xml:"part-name"`
}

type Part struct {
    XMLName xml.Name `xml:"part"`
    Id string `xml:"id,attr"`
    Measure []Measure `xml:"measure"`
}

type Measure struct {
    XMLName xml.Name `xml:"measure"`
    Number string `xml:"number,attr"`
    Attributes []Attribute  `xml:"attributes"`
    Note []Note `xml:"note"`
}

type Attribute struct {
    XMLName xml.Name `xml:"attributes"`
    Divisions string `xml:"divisions"`
    Key []Key `xml:"key"`
    Clef []Clef `xml:"clef"`
}

type Clef struct {
      XMLName xml.Name `xml:"clef"`
      Sign string `xml:"sign"`
      Line string `xml:"line"`
}

type Key struct {
    XMLName xml.Name `xml:"key"`
    Fifths string `xml:"fifths"`
}

type Note struct {
    XMLName xml.Name `xml:"note"`
    Chord string `xml:"chord,omitempty"`
    Pitch []Pitch `xml:"pitch,omitempty"`
    Rest string `xml:"rest,omitempty"`
    Duration string `xml:"duration"`
    Tie []Tie `xml:"tie,omitempty"`
    Type string `xml:"type"`
    Dot string `xml:"dot,omitempty"`
}

type Tie struct {
    XMLName xml.Name `xml:"tie"`
    Type string `xml:"type,attr"`
}

type Pitch struct {
    XMLName xml.Name `xml:"pitch"`
    Step string `xml:"step"`
    Alter string `xml:"alter,omitempty"`
    Octave string `xml:"octave"`
}

func main() {
    right_notes := []Note{
     
        {
            Pitch: []Pitch{
            {
                Step: "C",
                Alter: "0",
                Octave: "5",
            },
        },
            Duration: "1",
            Type: "quarter",
    },        
}
    left_notes := []Note{
        {
            Pitch: []Pitch{
            {
                Step: "E",
                Alter: "0",
                Octave: "4",

            },
        },
            Duration: "1",
            Type: "quarter",
    },
}

    notes := pitch_list()
    fmt.Println(notes)


    file, err := os.Open("data.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    right_next := "C_5" // right hand starting note
    right_flag := 0
    left_next := "E_4" // left hand starting note
    left_flag := 0

    lower_bound := indexOf("F_3", notes)
    upper_bound := indexOf("F_6", notes)

    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
    
        if len(line) > 1 {
            column_one := to_int(strings.Split(line,",")[0])
            column_two := to_int(strings.Split(line,",")[1])
            column_three := to_int(strings.Split(line,",")[2])
            column_four := to_int(strings.Split(line,",")[3])
            column_five := to_int(strings.Split(line,",")[4])
            column_six := to_int(strings.Split(line,",")[5])
            column_seven := to_int(alpha_num(strings.Split(line,",")[6]))

            right_next, right_flag = next_note(notes, right_next, column_one, column_two, lower_bound, upper_bound, right_flag)
       
            right_notes = note_append(right_notes,rh(right_next,column_three,rhythm(column_three)))

            left_next, left_flag = next_note(notes, left_next, column_five, column_six, lower_bound, upper_bound, left_flag)
   
            fmt.Println(column_one,column_two,column_three,column_four,column_five,column_six,column_seven)
      
            left_notes = note_append(left_notes, generate_left_note(left_next,column_four,column_seven,rhythm(column_seven),notes))
        }
    }

    score_partwise := ScorePartwise{
        Version: "4.0",
        PartList: []PartList{
            {
            ScorePart: []ScorePart{
                {
                    Id: "right",
                    PartName: "",
                },
                {
                    Id: "left",
                    PartName: "",

                },
            },
        },
    },
        Part: []Part{
            {
                Id: "right",
                Measure: []Measure{
                    {
                        Number: "1",
                        Attributes: []Attribute{
                            {
                                Divisions: "1",
                                Key: []Key{
                                    {
                                        Fifths: "0",
                                    },
                                },
                                Clef: []Clef{
                                    {
                                    Sign: "G",
                                    Line: "2",
                                    },
                                },
                            },
                        },
                        Note: right_notes,
                    },
                
            },
        },
             {
                Id: "left",
                Measure: []Measure{
                    {
                        Number: "1",
                        Attributes: []Attribute{
                            {
                                Divisions: "1",
                                Key: []Key{
                                    {
                                        Fifths: "0",
                                    },
                                },
                                Clef: []Clef{
                                    {
                                    Sign: "F",
                                    Line: "4",
                                    },
                                },
                            },
                        },
                        Note: left_notes,
                    },
                
            },
        },

    },
}


    xmlFile, err := os.Create("score.xml")
    if err != nil {
        fmt.Println("Error creating XML file: ", err)
        return
    }
    xmlFile.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>` + "\n" + `<!DOCTYPE score-partwise PUBLIC "-//Recordare/DTD MusicXML 4.0 Partwise//EN" "http://www.musicxml.org/dtds/partwise.dtd">` + "\n")
    encoder := xml.NewEncoder(xmlFile)
    encoder.Indent("", "\t")
    err = encoder.Encode(&score_partwise)
    if err != nil {
        fmt.Println("Error encoding XML to file: ", err)
        return
    }
    cmd := exec.Command("./script.sh")
    err = cmd.Run()
        if err != nil {
        log.Fatal(err)
    }
   
}
