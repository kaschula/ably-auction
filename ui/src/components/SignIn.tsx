import {FormEvent, MouseEvent, useState} from "react";

type SignUpProps = {
    setUserName: (name: string) => void
}

export function SignIn(props: SignUpProps) {
    const [userName, setUserName] = useState("")
    const submitHandler = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        props.setUserName(userName)
        setUserName("")
    }
    const onClickHandler = (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        props.setUserName(userName)
        setUserName("")
    }

    return (
        <div className="container-fluid">
            <div className={"row"}>
                <div className="col-md-12">
                    <h1>Ably Auction: Sign In</h1>
                </div>
            </div>
            <div className="mb-3 row">
                <form onSubmit={submitHandler}>
                    <div className="mb-3 row">
                        <label>Please enter you name: </label>
                        <div className={"col-md-10"}>
                            <input name={"userName"} onChange={(e) => setUserName(e.target.value)}/>
                        </div>
                    </div>
                    <div className="col-auto">
                        <button className={"btn btn-light"} disabled={userName === ""} type={"button"} onClick={onClickHandler}>Sign In</button>
                    </div>
                </form>
            </div>
        </div>
    )
}

