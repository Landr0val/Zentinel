// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title AuditLog
 * @dev Smart contract to store immutable audit logs of fraud alerts.
 * This contract serves as a "Black Box" for the Zentinel Fraud Detection System.
 */
contract AuditLog {
    address public owner;

    // Event emitted when a new alert is registered
    // indexed alertId allows filtering logs by specific alert ID
    event AlertRegistered(
        string indexed alertId,
        string dataHash,
        uint256 timestamp,
        address indexed recorder
    );

    modifier onlyOwner() {
        require(msg.sender == owner, "Caller is not the owner");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    /**
     * @dev Registers a new alert hash on the blockchain.
     * @param _alertId The unique identifier of the alert (UUID from the database).
     * @param _dataHash The SHA256 hash of the alert critical data (snapshot).
     */
    function registerAlert(string memory _alertId, string memory _dataHash) public onlyOwner {
        // We don't necessarily need to store data in state variables if we only need
        // the event logs for audit purposes. This saves significant gas.
        // The event log is permanent and immutable proof.
        emit AlertRegistered(_alertId, _dataHash, block.timestamp, msg.sender);
    }

    /**
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * Can only be called by the current owner.
     */
    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner is the zero address");
        owner = newOwner;
    }
}
