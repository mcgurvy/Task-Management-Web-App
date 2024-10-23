import React, { useState, useEffect } from 'react';

function App() {
  const [tasks, setTasks] = useState([]);
  const [title, setTitle] = useState('');

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    const response = await fetch('http://localhost:8080/tasks');
    const data = await response.json();
    setTasks(data);
  };

  const addTask = async () => {
    const response = await fetch('http://localhost:8080/task', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ title, done: false }),
    });

    if (response.ok) {
      fetchTasks();
      setTitle('');
    }
  };

  const toggleTask = async (id, done) => {
    await fetch(`http://localhost:8080/task/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ done: !done }),
    });

    fetchTasks();
  };

  const deleteTask = async (id) => {
    await fetch(`http://localhost:8080/task/${id}`, { method: 'DELETE' });
    fetchTasks();
  };

  return (
    <div className="App">
      <h1>Task Manager</h1>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="New Task"
      />
      <button onClick={addTask}>Add Task</button>
      <ul>
        {tasks.map((task) => (
          <li key={task.id}>
            <span
              style={{
                textDecoration: task.done ? 'line-through' : 'none',
                cursor: 'pointer',
              }}
              onClick={() => toggleTask(task.id, task.done)}
            >
              {task.title}
            </span>
            <button onClick={() => deleteTask(task.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
