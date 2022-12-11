# Day 11
Seemed fun to start with, didn't write a parser but just split strings on various lines. Didn't do the programming language concepts module and also not writing with a functional language so can't really do it easily.

I understood that as we no longer divided the level, we'd run into overflows soon so without much thought I switched to using the `Int` type provided by `math/big`

`type WorryLevel *big.Int` is completely different to `type WorryLevel = *big.Int` which took a while to debug. Got through the first few iterations okay to start with and then added a logging statement for the round number because it didn't seem to print round 1000 in the same way as it did before.

I've got lots of coursework deadlines to be doing so can't really spend ages on it today. I've looked at what is causing the issue with huge numbers and found that I need to somehow reduce the worry level to a manageable level, which doesn't affect the monkey that it gets passed to.

Initially I'm thinking to use GCD or something across all of the values, or maybe store the result of the mod if that would work.

Again, because I don't really have time to spend pondering this today, I need to get the LCM of all of the mod values and then just store the mod of each of the integers at each stage.

After about another hour of debugging, I've established that the inputs to my LCM function were completely wrong and I don't use the `DivisibleTrue` and `DivisibleFalse` values but the `DivisibleTest`... I think sometimes when I want to work quickly I don't take the time to work through it methodically.

Also I don't need the `math/big` library at all but due to not having any time, I won't refactor it out for this.

I should have realised that an integer overflow issue after 10 rounds on an input that will continue to grow with the round number is unlikely to be fixed by using a bigger integer and that I need to look at the problem more fundamentally.

All in all, I felt like part 1 was good and maintainable, then part 2 completely pulled the rug out and just became a mess.