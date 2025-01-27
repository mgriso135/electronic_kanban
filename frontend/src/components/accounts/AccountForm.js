import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const AccountForm = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [name, setName] = useState('');
    const [vatNumber, setVatNumber] = useState('');
    const [address, setAddress] = useState('');
    const [isEdit, setIsEdit] = useState(false);

    useEffect(() => {
        if (id) {
            setIsEdit(true);
            const fetchAccount = async () => {
              try {
                const response = await api.get(`/accounts/${id}`);
                setName(response.data.name);
                setVatNumber(response.data.vat_number);
                setAddress(response.data.address);
              } catch(error) {
                console.error('Error fetching account:', error)
              }
            };
            fetchAccount();
        }
      }, [id]);

      const handleSubmit = async (e) => {
        e.preventDefault();
        const accountData = {
          name: name,
          vat_number: vatNumber,
          address: address,
        };
        try {
          if (isEdit) {
            await api.put(`/accounts/${id}`, accountData);
          } else {
              await api.post('/accounts', accountData);
          }
            navigate('/accounts');
        } catch (error) {
          console.error('Error saving account:', error);
        }
      };


    return (
      <div>
        <h2>{isEdit ? 'Edit Account' : 'Create New Account'}</h2>
        <form onSubmit={handleSubmit}>
            <div>
              <label>Name:</label>
              <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
            </div>
            <div>
              <label>VAT Number:</label>
              <input type="text" value={vatNumber} onChange={(e) => setVatNumber(e.target.value)} />
            </div>
            <div>
                <label>Address:</label>
                <input type="text" value={address} onChange={(e) => setAddress(e.target.value)} />
              </div>
              <button type="submit">{isEdit ? 'Update Account' : 'Create Account'}</button>
              <Link to="/accounts">Cancel</Link>
          </form>
      </div>
    );
};

export default AccountForm;