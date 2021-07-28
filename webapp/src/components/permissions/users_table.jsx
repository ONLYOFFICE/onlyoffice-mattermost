/* eslint-disable react/prop-types */
/* eslint-disable no-console */
/* eslint-disable react/jsx-no-literals */
/* eslint-disable no-unused-vars */
/* eslint-disable max-nested-callbacks */
import React, {useState, useRef, useCallback} from 'react';
import {Table, Pane, Spinner, Button, TextInput} from 'evergreen-ui';

import {useSearchUsers} from 'utils/useSearchUsers';

// eslint-disable-next-line react/prop-types
export function UsersTable({fileInfoId, selected, setSelected}) {
    const [page, setPage] = useState(0);
    const [limit] = useState(25);

    const {
        data,
        hasMore,
        loading,
        error,
    } = useSearchUsers(page, limit, fileInfoId);

    const handleIntent = (name) => {
        if (selected === '*') {
            return 'success';
        }
        return selected.find((n) => n === name) ? 'success' : 'none';
    };

    const handleSelect = (name) => {
        if (selected === '*') {
            setSelected([name]);
        } else if (selected.find((n) => n === name)) {
            setSelected((prevSelected) => prevSelected.filter((item) => item !== name));
        } else {
            setSelected((prevSelected) => [...prevSelected, name]);
        }
    };

    const observer = useRef();
    const lastElementRef = useCallback((node) => {
        if (loading) {
            return;
        }

        if (observer.current) {
            observer.current.disconnect();
        }

        observer.current = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting && hasMore) {
                setPage((prevPage) => prevPage + 1);
            }
        });

        if (node) {
            observer.current.observe(node);
        }
    }, [loading, hasMore]);

    return (
        <div style={{position: 'relative', minWidth: '40%', height: '100%'}}>
            <Table style={{position: 'relative', width: '100%', height: '100%'}}>
                <Table.Head style={{display: 'flex', flexDirection: 'column', height: '15%', padding: 0}}>
                    <Table.HeaderCell style={{width: '100%', height: '50%', margin: '0', padding: '0'}}>
                        <TextInput
                            style={{width: '100%', height: '100%'}}
                            placeholder='Search by name'
                            disabled={loading || data.length === 0}
                        />
                    </Table.HeaderCell>
                    <Table.HeaderCell style={{width: '100%', height: '50%', margin: '0', padding: '0'}}>
                        <Button
                            style={{width: '100%', height: '100%'}}
                            onClick={() => (selected === '*' ? setSelected([]) : setSelected('*'))}
                            disabled={loading || data.length === 0}
                        >
                            Select ALL
                        </Button>
                    </Table.HeaderCell>
                </Table.Head>
                <Table.Body style={{position: 'relative', width: '100%', height: '85%'}}>
                    {loading && data.length === 0 ? (
                        <Pane
                            display='flex'
                            alignItems='center'
                            justifyContent='center'
                            height={120}
                        >
                            <Spinner/>
                        </Pane>
                    ) : (
                        <>
                            {data.map((name, index) => {
                                if (data.length === index + 1) {
                                    return (
                                        <Table.Row
                                            key={name}
                                            isSelectable={true}
                                            intent={handleIntent(name)}
                                            onSelect={() => handleSelect(name)}
                                        >
                                            <Table.TextCell ref={lastElementRef}>{name}</Table.TextCell>
                                        </Table.Row>
                                    );
                                }
                                return (
                                    <Table.Row
                                        key={name}
                                        isSelectable={true}
                                        intent={handleIntent(name)}
                                        onSelect={() => handleSelect(name)}
                                    >
                                        <Table.TextCell>{name}</Table.TextCell>
                                    </Table.Row>
                                );
                            },
                            )}
                            {loading && hasMore &&
                                <Table.Row
                                    display='flex'
                                    alignItems='center'
                                    justifyContent='center'
                                    height={120}
                                >
                                    <Spinner/>
                                </Table.Row>
                            }
                            {data.length === 0 &&
                                <div style={{display: 'flex', justifyContent: 'center'}}>
                                    <p>No users found</p>
                                </div>
                            }
                        </>
                    )}
                </Table.Body>
            </Table>
        </div >
    );
}
