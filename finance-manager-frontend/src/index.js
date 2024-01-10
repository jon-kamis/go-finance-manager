import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import App from './App';
import ErrorPage from './components/misc/ErrorPage';
import Home from './components/Home';
import Incomes from './components/income/Incomes';
import Users from './components/users/Users';
import Loans from './components/loans/Loans';
import Login from './components/auth/Login';
import Register from './components/auth/Register';
import ManageUser from './components/users/ManageUser';
import About from './components/misc/About';
import Bills from './components/bills/Bills';
import CreditCards from './components/creditcards/CreditCards';


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
        element: <Loans />,
      },
      {
        path: "/users/:userId/incomes",
        element: <Incomes />
      },
      {
        path: "/users/:userId/bills",
        element: <Bills />
      },
      {
        path: "/users/:userId/credit-cards",
        element: <CreditCards />
      }
    ]
  }
])

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);