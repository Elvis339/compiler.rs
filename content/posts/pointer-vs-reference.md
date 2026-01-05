+++
date = '2026-01-03T17:24:18+04:00'
draft = false
title = 'Familiar but Not Understood: Pointer vs Reference'
description = "Exploring concepts I thought I understood but didn't—starting with pointers vs references."
+++

I was watching [Tsoding](https://youtu.be/bzBU7iF2Rms?si=5pN4konbRqYlvUXS), where he questioned what a memory leak really means. He argues the definition is unclear for a few reasons:

1. "Not freed before exiting" doesn't help much because it is usually a good practice to pre-allocate a fixed size of memory, use it for the entirety of the program, and exiting doesn't cause any issue. So we don't really call this a leak.
2. "Unreachable memory" has limits too. This definition sounds better, but it breaks down in practice. In some languages, you can manipulate pointers directly, which can trick a garbage collector into thinking memory is unreachable when it's still in use. But even in languages like JavaScript where you don't have access to pointers, you can "leak" memory by forgetting to clean up callbacks that hold references to objects. Those objects should be unreachable, but they're not so memory keeps growing. You have a garbage collector and still leak memory. This is why Tsoding says there's no clear definition of a memory leak.

Tsoding suggests that a memory leak is more about what the programmer meant to do, not a strict rule. He also makes a deeper point that memory allocation itself is an artificial idea built on top of a fundamental computational model. So "leaks" are really just problems we created for ourselves.

> Funny comment from the chat: Memory leak is when I `Box::leak()`

This made me think about other concepts I assume I understand but actually don't. Right after watching, I looked at some Rust code and asked myself: what's really the difference between a reference and a pointer?

```rust
pub struct Holder<'a> {
    inner: &'a Vec<u8>,
    ptr: *const i32,
}
```

In university, I was taught they're basically the same "under the hood, references are just pointers." But if that's true, why does Rust have both references (`&T`) and raw pointers (`*const T`)? Why does C++ have both? And why do C and Go only have pointers?

## Semantics vs Implementation

When my professor said "references are just pointers under the hood," he wasn't wrong. At the assembly level, they often look identical:

```cpp
void by_ptr(int *p) { *p = 100; }
void by_ref(int &r) { r = 101; }
```

```asm
by_ptr(int*):
    mov QWORD PTR [rbp-8], rdi
    mov rax, QWORD PTR [rbp-8]
    mov DWORD PTR [rax], 100

by_ref(int&):
    mov QWORD PTR [rbp-8], rdi
    mov rax, QWORD PTR [rbp-8]
    mov DWORD PTR [rax], 101
```

Same instructions. But that's implementation, not semantics.

## Pointer Semantics

A pointer is a distinct value type. It has a defined size on a given platform (`sizeof(void*)`), and that size is the same regardless of what it points to. The value of a pointer is an actual memory address you can print it, compare it, reassign it.

Because pointers are true value types, they have their own rules independent of what they point to:

- Comparing pointers compares memory addresses, not the content at those addresses
- Copying a pointer copies the address, not the content
- Taking the address of a pointer gives you the address of the pointer itself, not the pointed-to content
- Dereferencing a pointer gives you access to the content

This is why you can have a null pointer. This is why you can reassign a pointer. The pointer exists as its own thing.

## Reference Semantics

A reference is not a distinct type. There is no object that "represents" a reference. You cannot take the address of the reference itself—if you try, you get the address of the bound instance. If you do `sizeof(T&)`, you get `sizeof(T)`. A reference is an alias for an existing object.

Because references are aliases, not value types:

- They must be bound to an instance at initialization—no null references
- They cannot be reassigned—assigning to a reference assigns to the bound instance
- Comparing references compares the referenced instances
- Taking the address gives you the address of the bound instance

The compiler is free to implement a reference however it wants, as long as it maintains these semantics. It can optimize a reference away completely. It can also implement it as a pointer internally. But a pointer must actually exist because you can take its address and reassign it.

```cpp
#include <iostream>

int main() {
    int a = 42;
    int b = 42;
    int* ptr1 = &a;
    int* ptr2 = &b;

    // Size is the same regardless of what it points to
    std::cout << "sizeof(int*): " << sizeof(int*) << "\n";
    std::cout << "sizeof(char*): " << sizeof(char*) << "\n";
    std::cout << "sizeof(double*): " << sizeof(double*) << "\n";
    // All print 8 on 64-bit

    // Comparing pointers compares addresses, not content
    std::cout << "a == b: " << (a == b) << "\n";           // true, same value
    std::cout << "ptr1 == ptr2: " << (ptr1 == ptr2) << "\n"; // false, different addresses

    // Pointer has its own address
    std::cout << "ptr1 value: " << ptr1 << "\n";      // address of a
    std::cout << "ptr1 address: " << &ptr1 << "\n";   // address of ptr1 itself indicating ptr is a value object stored somewhere

    // Can be null
    int* null_ptr = nullptr;
    std::cout << "null_ptr: " << null_ptr << "\n";    // 0

    // Can be reassigned
    ptr1 = &b;
    std::cout << "ptr1 after reassign: " << *ptr1 << "\n"; // 42 (b's value)

    return 0;
}
```

## References Across Languages

Not all references are the same. C++ references are mutable you can change the value through them. Rust references are read-only by default (`&T`), and you need `&mut T` for mutation. C and Go don't have references at all, just pointers.

| | C | Go | C++ | Rust |
|---|---|---|---|---|
| Pointers | `*T` | `*T` | `T*` | `*const T`, `*mut T` |
| References | - | - | `T&` (mutable) | `&T` (read-only), `&mut T` |

So when someone says "reference," ask which language. The semantics are different.

## The "Reference Semantics" Trap

People say both pointers and references have "reference semantics." This is sloppy. A pointer doesn't have reference semantics a pointer is a value type that *models* reference semantics through `*` and `->`.

Trick question: does Python have reference or value semantics?

```python
def reassign(x):
    x = [99]

def mutate(x):
    x.append(99)

a = [1, 2, 3]
reassign(a)
print(a)  # [1, 2, 3] — unchanged

b = [1, 2, 3]
mutate(b)
print(b)  # [1, 2, 3, 99] — changed
```

Answer: Python has value semantics, but the only value type you can assign is a pointer to an object. In `reassign`, you're reassigning the local pointer `x` to a new list. In `mutate`, you're dereferencing the pointer and changing the object it points to.

## Why Does This Matter?

In C++ reference kind of means indirect access of a pointer but with better properties for the common case: aliasing a single object. References have nicer syntax (no `*` and `->`), no null pointer edge cases, and no accidental reassignment.

In C++ reference is "a pointer with restrictions." But the underlying implentation is just the aliasing properties, not anything about the actual type pointer. That's the gap between implementation and semantics.

## Conclusion

Tsoding's point about memory leaks applies here too. Memory allocation is an artificial concept we built on top of how computers work and now we have to deal with "leaks" as a consequence. The same is true for pointers and references. We created these abstractions, and now we deal with the confusion they might cause.
"Under the hood, references are just pointers" is true at the assembly level. But it's a surface level understanding. Pointers and references have different semantics different rules about nullability, reassignment, and identity. And those semantics change between languages. A reference in C++ is mutable. A reference in Rust is read-only by default. C and Go don't have references at all.

If you move between languages and assume "reference" means the same thing everywhere, it can bite you.

### The mental model

Everything is pass by value. Always.

When you pass something:

1. Primitive types copy the value. Two independent values.
2. Pointer/reference to heap object copy the address. Two independent pointers to same data.
   - Dereference and mutate (`p.val = 5`, `obj.x = 5`) → affects shared data
   - Reassign (`p = new`, `obj = new`) → only affects your local copy

C++ has an exception

- `void foo(T x)` → copy (same as above)
- `void foo(T* p)` → copy of pointer (same as above)
- `void foo(T& r)` → alias, no copy, IS the caller's variable

**The flowchart:**

```
Passing argument to function
            |
            v
   Is it C++ reference (&T)?
           /    \
         yes     no
          |       |
          v       v
      no copy    copy value into new local variable
      (alias)               |
                            v
                  is it a pointer?
                      /       \
                    yes        no
                     |          |
                     v          v
             two pointers    two independent
             same data       values
```
