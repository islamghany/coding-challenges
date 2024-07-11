# Takewayes

## Unix Philosophy

**Ken Thompson**:

- Small is beautiful.
- Write programs that do one thing and do it well.
- Write programs to work together.
- Write programs to handle text streams, because that is a universal interface.

**Rob Pike**:

- Rule 1: you can't tell where a program is going to spend its time. Bottlenecks occur in surprising places, so don't try to second guess and put in a speed hack until you've proven that's where the bottleneck is.
- Rule 2: Measure. Don't tune for speed until you've measured, and even then don't unless one part of the code overwhelms the rest.
- Rule 3: Fancy algorithms are slow when n is small, and n is usually small. Fancy algorithms have big constants. Until you know that n is frequently going to be big, don't get fancy. (Even if n does get big, use Rule 2 first.)
- Rule 4: Fancy algorithms are buggier than simple ones, and they're much harder to implement. Use simple algorithms as well as simple data structures.
- Rule 5: Data dominates. If you've chosen the right data structures and organized things well, the algorithms will almost always be self-evident. Data structures, not algorithms, are central to programming.
- Rule 6: There is no Rule 6.

### Rule of Modularity: Write simple parts connected by clean interfaces.

we can use the Unix philosophy: Write programs that do one thing and do it well. Write programs to work together. Write programs to handle text streams, because that is a universal interface.

### Rule of Clarity: Clarity is better than cleverness.

because maintenance is the ultimate goal of all code, and clever code is harder to maintain.

### Rule of Composition: Desigggn programs to be connected to other programs.

To make programs composable, make them independent. A program on one end of a text stream should care as little as possible about the program on the other end. It should be made easy to replace one end with a completely different implementation without disturbing the other.

### Rule of Sepration: Seprate policy from mechanism; separate interfaces from engines.

The engine is the core of the program. It implements the program's functionality. The interface is just the way the engine is called. The engine should be designed to be testable in isolation. The interface should be as thin as possible.

### Rule of Simplicity: Design for simplicity; add complexity only where you must.

the simplest solution is usually the best. If you can make your program simpler, you should. Simplicity is the most important consideration in a design.
and complexity is the worst enemy of security and reliability.

### Rule of Parsimony: Write a big program only when it is clear by demonstration that nothing else will do.

programs should be small, but the Unix philosophy does not say that they should be tiny. If you need a big program, you should write one. But don't write a big one when a small one will do.

### Rule of Transparency: Design for visibility to make inspection and debugging easier.

The Unix philosophy is to make programs as transparent as possible. If you can, make the internals of your program visible to the outside world. If you can't do that, make them visible to the debugger.

The objective of designing for transparency and discoverability should also encourage simple interfaces that can easily be manipulated by other programs — in particular, test and monitoring harnesses and debugging scripts.

### Rule of Robustness: Robustness is the child of transparency and simplicity.

Software is said to be robust when it performs well under unexpected conditions which stress the designer's assumptions, as well as under normal conditions.

### Rule of Representation: Fold knowledge into data, so program logic can be stupid and robust.

Data is more tractable than program logic. It follows that where you see a choice between complexity in data structures and complexity in code, choose the former. More: in evolving a design, you should actively seek ways to shift complexity from code to data.

Encode as much information as possible into data structures rather than complex logic.

### Rule of Least Surprise: In interface design, always do the least surprising thing.

Design interfaces that behave in a way users expect.

The Principle of Least Surprise (or consistency principle) is the idea that a user shouldn't be surprised by the way an interaction or object works in an interface or design. This means prioritizing functionality and use over things like consistency to avoid astonishing or surprising your user.

### Rule of Silence: When a program has nothing surprising to say, it should say nothing.

Avoid unnecessary output.

When a program has nothing surprising to say, it should say nothing. The Unix philosophy is to write programs that handle error conditions, so that the user does not have to worry about them.

### Rule of Repair: When you must fail, fail noisily and as soon as possible.

Ensure errors are detected and reported early.

When you must fail, fail noisily and as soon as possible. This is the only way to get the user's attention, and to let the user know that the problem is serious and should be attended to.

### Rule of Economy: Programmer time is expensive; conserve it in preference to machine time.

Optimize for developer productivity, even if it means using more computational resources.

Programmer time is expensive; conserve it in preference to machine time. The Unix philosophy is to design programs that make the easy things easy and the hard things possible. If you can make the easy things easy, you should.

### Rule of Generation: Avoid hand-hacking; write programs to write programs when you can.

Automate code generation and other repetitive tasks.

Avoid hand-hacking; write programs to write programs when you can. This rule is the natural consequence of the Rule of Economy. It has long been a Unix tradition to avoid interactive programming with a text editor. The idea is to write a program that writes the program you want, and then to execute the program that writes the program.

### Rule of Optimization: Prototype before polishing. Get it working before you optimize it.

Make it work before you make it fast.

“Premature optimization is the root of all evil”

Prototype before polishing. Get it working before you optimize it. The idea is to do the simplest thing that could possibly work, and then to improve it if necessary.

### Rule of Diversity: Distrust all claims for “one true way”.

Be open to different approaches and solutions.

the Unix tradition includes a healthy mistrust of “one true way” approaches to software design or implementation. It embraces multiple languages, open extensible systems, and customization hooks everywhere.

### Rule of Extensibility: Design for the future, because it will be here sooner than you think.

Design for change.
