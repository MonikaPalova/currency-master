document.addEventListener("load", _loadNav());

function _loadNav() {
    document.body.insertAdjacentHTML('afterbegin', `
                                                    <nav id="nav-bar" class="column">
                                                        <div class="nav-link">
                                                            <a href="">Profile</a>
                                                        </div> 
                                                        <div class="nav-link">
                                                            <a href="">Wallet</a>
                                                        </div> 
                                                        <div class="nav-link">
                                                            <a href="">Current Prices</a>
                                                        </div> 
                                                    </nav>
`);
}