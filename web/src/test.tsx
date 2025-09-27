// i want to write a react function that export hello world in frontend

import React, { useCallback } from "react";
import { useState, useEffect, useRef, useMemo } from "react";


export default function Test() {
  const [count, setCount] = useState(0);

  const onClick = () => {
    setCount(count + 1); // captures 'count' from that render
    console.log(count); // will log old 'count' value
  }

//   const onClick = useCallback(() => {
//     setCount(c => c + 1); 
//     console.log(count); 
//   }
// , []); // no dependencies, created once

  // ...after some renders, onClick may still add 1 to an old 'count'
  return <button onClick={onClick}>Click me</button>;
}

export function TestuseEffect() {
  const [users, setUsers] = useState<Array<{ id: number; email: string }>>([]);

  useEffect(() => {
    const ctrl = new AbortController();

    (async () => {
      try {
        const res = await fetch("http://localhost:8080/api/user/auth/all", { signal: ctrl.signal });
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }
        const data = await res.json();
        setUsers(data);
      } catch (err) {
        if (err instanceof Error) {
          if (err.name === "AbortError") {
            console.log("Fetch aborted");
          } else {
            console.error("Failed to fetch users:", err.message);
          }
        } else {
          console.error("An unknown error occurred:", err);
        }
      }
    })();

    return () => ctrl.abort(); // Cleanup on unmount or dependency change
  }, []); // [] => run once


  return <ul>{users.map(u => <li key={u.id}>{u.email}</li>)}</ul>;
}

export function Counter() {
  const [count, setCount] = useState(0); // Initialize state with 0

  return (
    <div>
      <h1>Counter: {count}</h1>
      <button onClick={() => setCount(count + 1)}>Increment</button>
      <button onClick={() => setCount(count - 1)}>Decrement</button>
      <button onClick={() => setCount(0)}>Reset</button>
    </div>
  );
}

export function Card() {
  return (
    <div className="card">
      <h2>Card Title</h2>
      <p>This is a card component.</p>
    </div>
  );
}

export function AutoFocusInput() {
  // const inputRef = useRef<HTMLInputElement | null>(null);
  const inputRef = "true";

  // useEffect(() => {
  //   inputRef.current?.focus();
  // }, []);

  return <input value={inputRef} placeholder="I get focus on mount" />;
}

function Timer() {
  const countRef = useRef(0);

  function increment() {
    countRef.current += 1;
    console.log("Count:", countRef.current);
  }

  return <button onClick={increment}>Add</button>;
}

const fibCache: { [key: number]: number } = {};

function slowFib(n: number): number {
  if (n in fibCache) {
    console.log(`Returning cached result for n = ${n}`);
    return fibCache[n];
  }

  console.log(`Calculating slowFib for n = ${n}`);
  fibCache[n] = n <= 1 ? n : slowFib(n - 1) + slowFib(n - 2);
  return fibCache[n];
}

export function FibDemo() {
  const [n, setN] = useState(35);

  const value = useMemo(() => slowFib(n), [n]);

  return (
    <>
      <input type="number" value={n} onChange={e => setN(+e.target.value)} />
      <div>fib({n}) = {value}</div>
    </>
  );
}


export function WithoutUseRef() {
  let ticks = 0;            // resets to 0 on every render
  function addTick() { ticks += 1; console.log(ticks); }
  return <button onClick={addTick}>Add tick {ticks}</button>;
}

export function WithUseRef() {
  const ticksRef = React.useRef(0);

  function addTick() {
    ticksRef.current += 1;       // updates value
    console.log(ticksRef.current); // UI doesn't re-render
  }

  return <button onClick={addTick}>Add tick (check console){ticksRef.current}</button>;
}

type State = { name: string; age: number; submitted: boolean };
type Action =
  | { type: "SET_NAME", payload: string }
  | { type: "SET_AGE", payload: number }
  | { type: "SUBMIT" };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case "SET_NAME": return { ...state, name: action.payload };
    case "SET_AGE": return { ...state, age: action.payload };
    case "SUBMIT": return { ...state, submitted: true };
    default: return state;
  }
}

export function FormWithReducer() {
  const [state, dispatch] = React.useReducer(reducer, {
    name: "",
    age: 0,
    submitted: false
  });

  return (
    <>
      {!state.submitted ? (
        <>
          <input
            placeholder="Name"
            value={state.name}
            onChange={(e) =>
              dispatch({ type: "SET_NAME", payload: e.target.value })
            }
          />
          <input
            placeholder="Age"
            type="number"
            value={state.age}
            onChange={(e) =>
              dispatch({ type: "SET_AGE", payload: Number(e.target.value) })
            }
          />
          <button onClick={() => dispatch({ type: "SUBMIT" })}>Submit</button>
        </>
      ) : (
        <div>
          <h2>Submitted Data:</h2>
          <p>Name: {state.name}</p>
          <p>Age: {state.age}</p>
        </div>
      )}
    </>
  );
}