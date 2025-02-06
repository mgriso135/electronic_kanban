import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

const AccountList = () => {
    const [accounts, setAccounts] = useState(null); //set to null initially

    useEffect(() => {
        const fetchAccounts = async () => {
            try {
                const response = await api.get('/accounts');
                console.log("API response:", response);
                setAccounts(response.data);
            } catch (error) {
                console.error('Error fetching accounts:', error);
                setAccounts([]); // set accounts to empty array in case of error
            }
        };

        fetchAccounts();
    }, []);

    const handleDelete = async (id) => {
      try {
          await api.delete(`/accounts/${id}`);
          setAccounts(accounts.filter(account => account.id !== id));
      } catch (error) {
          console.error('Error deleting account:', error);
      }
    };

    return (
      <div>
          <h2>Accounts</h2>
          <Link to="/accounts/new">Create New Account</Link>
            {accounts && accounts.length > 0 ? (
                <table>
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>VAT Number</th>
                        <th>Address</th>
                        <th>Actions</th>
                    </tr>
                    </thead>
                    <tbody>
                    {accounts.map(account => (
                        <tr key={account.id}>
                            <td>{account.id}</td>
                            <td>{account.name}</td>
                            <td>{account.vat_number}</td>
                            <td>{account.address}</td>
                            <td>
                              <Link to={`/accounts/${account.id}/edit`}>Edit</Link>
                              <button onClick={() => handleDelete(account.id)}>Delete</button>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                  </table>
                ) : (
                  <p>No accounts found</p>
            )}
        </div>
    );
};

export default AccountList;