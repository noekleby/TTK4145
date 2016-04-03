#pragma once

// Returns 0 on init failure
int ioInit(void);

void ioSetBit(int channel);
void ioClearBit(int channel);

int ioReadBit(int channel);

int ioReadAnalog(int channel);
void ioWriteAnalog(int channel, int value);

