package main

import (
	"errors"
	"fmt"
)

func main() {
	var (
		x_expLv, uLow, uMedium, uHigh    float32 //exploration level
		x_expTime, uFast, uNormal, uSlow float32 //exploration time
		poor, average, awesome           float32 //conclusion
		reward                           float32 //reward outcome
	)
	x_expLv = 2.5
	x_expTime = 1.0
	fmt.Println("===\tinput\t===")
	fmt.Println("x_expLv: ", x_expLv)
	fmt.Println("x_expTime: ", x_expTime)

	//get fuzzification for exploration level and time
	uLow, uMedium, uHigh = fuzzification(x_expLv, 40.0, 60.0, 80.0)
	uFast, uNormal, uSlow = fuzzification(x_expTime, 15.0, 30.0, 45.0)

	fmt.Println("===\tfuzzification of exploration level\t===")
	fmt.Printf("low:%2f\n", uLow)
	fmt.Printf("Medium:%2f\n", uMedium)
	fmt.Printf("High:%2f\n", uHigh)

	fmt.Println("===\tfuzzification of exoploration time\t===")
	fmt.Printf("Fast:%2f\n", uFast)
	fmt.Printf("Normal:%2f\n", uNormal)
	fmt.Printf("Slow:%2f\n", uSlow)

	// and rules
	poor, average, awesome = 0, 0, 0 // default value for
	andRules(uLow, uFast, &poor)     // low && fast -> poor
	andRules(uLow, uNormal, &poor)   // low && normal -> poor

	andRules(uMedium, uSlow, &average)
	andRules(uMedium, uFast, &average)   // medium && fast -> average
	andRules(uMedium, uNormal, &average) // medium && normal -> average

	andRules(uMedium, uSlow, &awesome) // medium && slow -> awesome

	andRules(uHigh, uFast, &average)   // high && fast -> average
	andRules(uHigh, uNormal, &awesome) // high && normal -> awesome
	andRules(uHigh, uSlow, &awesome)   // high && slow -> awesome

	fmt.Println("===\tconclusion\t===")
	fmt.Println("poor: ", poor)
	fmt.Println("average: ", average)
	fmt.Println("awesome: ", awesome)

	// determine sample
	// 40.0, 60.0, 80.0
	smp1, smp2, smp3 := determine_sample(0.0, 40.0, 80.0, 100.0)
	fmt.Println("===\tSample value for defuzzification\t===")
	fmt.Println("smp1: ", smp1)
	fmt.Println("smp2: ", smp2)
	fmt.Println("smp3: ", smp3)

	// defuzzification
	reward, err := defuzzification(poor, average, awesome, smp1, smp2, smp3)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("reward: ", reward)
}

func fuzzification(x float32, a float32, b float32, c float32) (u1, u2, u3 float32) {
	var uLeft float32   //membership value for half trapezoid (left)
	var uCenter float32 //membership value for trapezoid
	var uRight float32  //membership value for half trapezoid (right)

	var (
		left_c, left_d      float32 // left trapeziod
		ctr_a, ctr_b, ctr_c float32 // center triangle
		rigth_a, right_b    float32 //right trapezoid
	)

	left_c = a
	left_d = b
	ctr_a = a
	ctr_b = b
	ctr_c = c
	rigth_a = b
	right_b = c

	//left trapezoid
	if x <= left_c {
		uLeft = 1
	} else if x > left_c && x < left_d {
		uLeft = (left_d - x) / (left_d - left_c)
	} else {
		uLeft = 0
	}
	fmt.Println("u_left: ", uLeft)

	// center triangle
	if (x <= ctr_a) || (x >= ctr_c) {
		uCenter = 0
	} else if (x > ctr_a) && (x < ctr_b) {
		// upside function
		uCenter = (x - ctr_a) / (ctr_b - ctr_a)
	} else if (x > ctr_b) && (x < ctr_c) {
		// down function
		uCenter = (ctr_c - x) / (ctr_c - ctr_b)
	} else if x == ctr_b {
		uCenter = 1
	}
	fmt.Println("u-center: ", uCenter)

	//right trapezoid
	if x <= rigth_a {
		uRight = 0
	} else if x > rigth_a && x < right_b {
		uRight = (x - rigth_a) / (right_b - rigth_a)
	} else {
		uRight = 1
	}
	fmt.Println("u_right: ", uRight)

	u1 = uLeft
	u2 = uCenter
	u3 = uRight
	return u1, u2, u3
}

func andRules(f1 float32, f2 float32, c *float32) {
	var tmp float32 = 0
	if f1 >= f2 {
		tmp = f2
	} else {
		tmp = f1
	}

	if tmp > *c {
		*c = tmp
	}

}

func determine_sample(a, b, c, d float32) (s1, s2, s3 float32) {
	s1 = (a + b) / 2
	s2 = (b + c) / 2
	s3 = (c + d) / 2
	return s1, s2, s3
}

func defuzzification(u1, u2, u3, s1, s2, s3 float32) (result float32, err error) {

	e1, e2 := nilException(u1, u2, u3), nilException(s1, s2, s3)
	if (e1 != nil) || (e2 != nil) {
		return 0.0, e1
	}

	result = ((s1 * u1) + (s2 * u2) + (s3 * u3)) / (u1 + u2 + u3)
	return result, nil
}

func nilException(f ...float32) error {
	var tmp float32
	for _, v := range f {
		tmp += v
	}
	if tmp == 0 {
		return errors.New("input cannot be zero")
	}
	return nil
}
