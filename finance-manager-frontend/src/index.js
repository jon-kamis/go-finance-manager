import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import App from './App';
import ManageLoan from './components/ManageLoan';
import ErrorPage from './components/ErrorPage';
import Home from './components/Home';
import Users from './components/Users';
import Loans from './components/Loans';
import Login from './components/Login';
import NewLoan from './components/NewLoan';
import Register from './components/Register';
import ManageUser from './components/ManageUser';
import About from './components/About';

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <Home /> },
      {
        path: "/about",
        element: <About />
      },
      {
        path: "/users",
        element: <Users />
      },
      {
        path: "admin/users/:id",
        element: <ManageUser />
      },
      {
        path: "/login",
        element: <Login />
      },
      {
        path: "/register",
        element: <Register />
      },
      {
        path: "/users/:userId/loans",
        element: <Loans />
      },
      {
        path: "/users/:userId/loans/new",
        element: <NewLoan />
      },
      {
        path: "/users/:userId/loans/:loanId",
        element: <ManageLoan />
      }
    ]
  }
])

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
 <React.StrictMode>
    <RouterProvider router={router}/>
  </React.StrictMode>
);